package service

import (
	"encoding/json"
	"fmt"

	"blog_backend/internal/model/dto/request"
	"blog_backend/internal/model/dto/response"
	"blog_backend/internal/repository"
)

// 关于页面设置键
const (
	SettingKeyAboutProfile    = "about_profile"
	SettingKeyAboutHighlights = "about_highlights"
	SettingKeyAboutSkills     = "about_skills"
	SettingKeyAboutTimeline   = "about_timeline"
	SettingKeyAboutPhilosophy = "about_philosophy"
)

// systemSettingService 系统设置服务实现
type systemSettingService struct {
	settingRepo repository.SystemSettingRepository
}

// NewSystemSettingService 创建系统设置服务
func NewSystemSettingService(settingRepo repository.SystemSettingRepository) SystemSettingService {
	return &systemSettingService{
		settingRepo: settingRepo,
	}
}

// GetAboutPage 获取关于页面数据
func (s *systemSettingService) GetAboutPage() (*response.AboutPageResponse, error) {
	// 获取所有关于页面相关的设置
	keys := []string{
		SettingKeyAboutProfile,
		SettingKeyAboutHighlights,
		SettingKeyAboutSkills,
		SettingKeyAboutTimeline,
		SettingKeyAboutPhilosophy,
	}

	settings, err := s.settingRepo.FindByKeys(keys)
	if err != nil {
		return nil, fmt.Errorf("获取关于页面数据失败: %w", err)
	}

	// 转换为 map 便于处理
	settingsMap := make(map[string]string)
	for _, setting := range settings {
		settingsMap[setting.Key] = setting.Value
	}

	// 构建响应
	result := &response.AboutPageResponse{}

	// 解析个人信息
	if profileStr, ok := settingsMap[SettingKeyAboutProfile]; ok {
		var profile response.AboutProfile
		if err := json.Unmarshal([]byte(profileStr), &profile); err == nil {
			result.Profile = &profile
		}
	}

	// 解析特点
	if highlightsStr, ok := settingsMap[SettingKeyAboutHighlights]; ok {
		var highlights []string
		if err := json.Unmarshal([]byte(highlightsStr), &highlights); err == nil {
			result.Highlights = highlights
		}
	}

	// 解析技术栈
	if skillsStr, ok := settingsMap[SettingKeyAboutSkills]; ok {
		var skills []*response.AboutSkill
		if err := json.Unmarshal([]byte(skillsStr), &skills); err == nil {
			result.Skills = skills
		}
	}

	// 解析经历时间轴
	if timelineStr, ok := settingsMap[SettingKeyAboutTimeline]; ok {
		var timeline []*response.AboutTimeline
		if err := json.Unmarshal([]byte(timelineStr), &timeline); err == nil {
			result.Timeline = timeline
		}
	}

	// 解析写作理念
	if philosophyStr, ok := settingsMap[SettingKeyAboutPhilosophy]; ok {
		var philosophy response.AboutPhilosophy
		if err := json.Unmarshal([]byte(philosophyStr), &philosophy); err == nil {
			result.Philosophy = &philosophy
		}
	}

	return result, nil
}

// UpdateAboutPage 更新关于页面数据
func (s *systemSettingService) UpdateAboutPage(req *request.UpdateAboutPageRequest) error {
	settings := make(map[string]string)

	// 序列化个人信息
	if req.Profile != nil {
		profileBytes, err := json.Marshal(req.Profile)
		if err != nil {
			return fmt.Errorf("序列化个人信息失败: %w", err)
		}
		settings[SettingKeyAboutProfile] = string(profileBytes)
	}

	// 序列化特点
	if req.Highlights != nil {
		highlightsBytes, err := json.Marshal(req.Highlights)
		if err != nil {
			return fmt.Errorf("序列化特点失败: %w", err)
		}
		settings[SettingKeyAboutHighlights] = string(highlightsBytes)
	}

	// 序列化技术栈
	if req.Skills != nil {
		skillsBytes, err := json.Marshal(req.Skills)
		if err != nil {
			return fmt.Errorf("序列化技术栈失败: %w", err)
		}
		settings[SettingKeyAboutSkills] = string(skillsBytes)
	}

	// 序列化经历时间轴
	if req.Timeline != nil {
		timelineBytes, err := json.Marshal(req.Timeline)
		if err != nil {
			return fmt.Errorf("序列化经历时间轴失败: %w", err)
		}
		settings[SettingKeyAboutTimeline] = string(timelineBytes)
	}

	// 序列化写作理念
	if req.Philosophy != nil {
		philosophyBytes, err := json.Marshal(req.Philosophy)
		if err != nil {
			return fmt.Errorf("序列化写作理念失败: %w", err)
		}
		settings[SettingKeyAboutPhilosophy] = string(philosophyBytes)
	}

	// 批量保存
	for key, value := range settings {
		if err := s.settingRepo.Set(key, value); err != nil {
			return fmt.Errorf("保存设置 %s 失败: %w", key, err)
		}
	}

	return nil
}

// GetSetting 获取单个设置
func (s *systemSettingService) GetSetting(key string) (string, error) {
	setting, err := s.settingRepo.FindByKey(key)
	if err != nil {
		return "", fmt.Errorf("获取设置 %s 失败: %w", key, err)
	}
	if setting == nil {
		return "", nil
	}
	return setting.Value, nil
}

// SetSetting 设置单个值
func (s *systemSettingService) SetSetting(key, value string) error {
	return s.settingRepo.Set(key, value)
}

// GetAllSettings 获取所有设置
func (s *systemSettingService) GetAllSettings() ([]*response.SettingResponse, error) {
	settings, err := s.settingRepo.FindAll()
	if err != nil {
		return nil, fmt.Errorf("获取所有设置失败: %w", err)
	}

	result := make([]*response.SettingResponse, 0, len(settings))
	for _, setting := range settings {
		result = append(result, &response.SettingResponse{
			Key:   setting.Key,
			Value: setting.Value,
		})
	}

	return result, nil
}
