package service

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"sync"

	"xx-backend/internal/model"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type UserService struct {
	db           *gorm.DB
	redis        *redis.Client
	kafkaService *KafkaService
	mu           sync.RWMutex
}

func NewUserService(db *gorm.DB, redis *redis.Client, kafkaService *KafkaService) *UserService {
	return &UserService{
		db:           db,
		redis:        redis,
		kafkaService: kafkaService,
	}
}

// GetUsers 获取用户列表（支持分页和搜索）
func (s *UserService) GetUsers(page, pageSize int, search string) ([]model.User, int64, error) {
	var users []model.User
	var total int64

	query := s.db.Model(&model.User{}).Preload("Role")

	if search != "" {
		query = query.Where("username LIKE ? OR nickname LIKE ?", "%"+search+"%", "%"+search+"%")
	}

	// 获取总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUser 根据ID获取用户
func (s *UserService) GetUser(id int) (*model.User, error) {
	var user model.User
	if err := s.db.Preload("Role").First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

// CreateUser 创建用户
func (s *UserService) CreateUser(user *model.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查用户名是否已存在
	var count int64
	if err := s.db.Model(&model.User{}).Where("username = ?", user.Username).Count(&count).Error; err != nil {
		return err
	}
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}

	err := s.db.Create(user).Error
	if err != nil {
		return err
	}

	// 记录用户注册事件到Kafka
	if s.kafkaService != nil {
		if err := s.kafkaService.LogUserRegister(user.ID, user.Username, user.Email); err != nil {
			// 记录Kafka错误但不影响用户创建流程
			fmt.Printf("Failed to log user register to Kafka: %v\n", err)
		}
	}

	return nil
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(id int, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取用户信息用于Kafka记录
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	err := s.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
	if err != nil {
		return err
	}

	// 记录用户更新事件到Kafka
	if s.kafkaService != nil {
		if err := s.kafkaService.LogUserUpdate(user.ID, user.Username, updates); err != nil {
			// 记录Kafka错误但不影响用户更新流程
			fmt.Printf("Failed to log user update to Kafka: %v\n", err)
		}
	}

	return nil
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 获取用户信息用于Kafka记录
	var user model.User
	if err := s.db.First(&user, id).Error; err != nil {
		return err
	}

	err := s.db.Delete(&model.User{}, id).Error
	if err != nil {
		return err
	}

	// 记录用户删除事件到Kafka
	if s.kafkaService != nil {
		data := map[string]interface{}{
			"user_id":  user.ID,
			"username": user.Username,
			"action":   "delete",
		}
		if err := s.kafkaService.client.SendUserEvent("user_delete", data); err != nil {
			// 记录Kafka错误但不影响用户删除流程
			fmt.Printf("Failed to log user delete to Kafka: %v\n", err)
		}
	}

	return nil
}

// GetProfile 获取用户资料
func (s *UserService) GetProfile(userID int) (*model.User, error) {
	var user model.User
	err := s.db.Preload("Role").First(&user, userID).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// GetRoles 获取角色列表
func (s *UserService) GetRoles() ([]model.Role, error) {
	var roles []model.Role
	if err := s.db.Find(&roles).Error; err != nil {
		return nil, err
	}
	return roles, nil
}

// CreateRole 创建角色
func (s *UserService) CreateRole(role *model.Role) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Create(role).Error
}

// UpdateRole 更新角色
func (s *UserService) UpdateRole(id uint, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Model(&model.Role{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteRole 删除角色
func (s *UserService) DeleteRole(id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Delete(&model.Role{}, id).Error
}

// GetMenus 获取菜单列表
func (s *UserService) GetMenus() ([]model.Menu, error) {
	var menus []model.Menu
	if err := s.db.Order("sort").Find(&menus).Error; err != nil {
		return nil, err
	}
	return menus, nil
}

// CreateMenu 创建菜单
func (s *UserService) CreateMenu(menu *model.Menu) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Create(menu).Error
}

// UpdateMenu 更新菜单
func (s *UserService) UpdateMenu(id uint, updates map[string]interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Model(&model.Menu{}).Where("id = ?", id).Updates(updates).Error
}

// DeleteMenu 删除菜单
func (s *UserService) DeleteMenu(id uint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.db.Delete(&model.Menu{}, id).Error
}

// BatchProcessUsers 批量处理用户（多线程示例）
func (s *UserService) BatchProcessUsers(userIDs []int, processor func(*model.User) error) error {
	var wg sync.WaitGroup
	errChan := make(chan error, len(userIDs))

	// 创建工作池
	workerCount := 5
	userChan := make(chan int, len(userIDs))

	// 启动工作协程
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for userID := range userChan {
				user, err := s.GetUser(userID)
				if err != nil {
					errChan <- err
					continue
				}
				if err := processor(user); err != nil {
					errChan <- err
				}
			}
		}()
	}

	// 发送任务
	for _, userID := range userIDs {
		userChan <- userID
	}
	close(userChan)

	// 等待所有工作完成
	wg.Wait()
	close(errChan)

	// 检查是否有错误
	for err := range errChan {
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) Register(username, password, email string) error {
	// 检查用户名是否已存在
	var count int64
	s.db.Model(&model.User{}).Where("username = ?", username).Count(&count)
	if count > 0 {
		return fmt.Errorf("用户名已存在")
	}
	// 密码加密
	hash := md5.Sum([]byte(password))
	user := model.User{
		Username: username,
		Password: hex.EncodeToString(hash[:]),
		Email:    email,
		Status:   1,
		RoleID:   2, // 普通用户
	}

	err := s.db.Create(&user).Error
	if err != nil {
		return err
	}

	// 记录用户注册事件到Kafka
	if s.kafkaService != nil {
		if err := s.kafkaService.LogUserRegister(user.ID, user.Username, user.Email); err != nil {
			// 记录Kafka错误但不影响注册流程
			fmt.Printf("Failed to log user register to Kafka: %v\n", err)
		}
	}

	return nil
}
