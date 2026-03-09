//go:build wireinject

package di

import (
	aiApi "github.com/GoSimplicity/CloudOps/internal/ai/api"
	aiDao "github.com/GoSimplicity/CloudOps/internal/ai/dao"
	aiService "github.com/GoSimplicity/CloudOps/internal/ai/service"
	"github.com/GoSimplicity/CloudOps/internal/ai/vector"
	k8sApi "github.com/GoSimplicity/CloudOps/internal/k8s/api"
	k8sClient "github.com/GoSimplicity/CloudOps/internal/k8s/client"
	k8sDao "github.com/GoSimplicity/CloudOps/internal/k8s/dao"
	k8sService "github.com/GoSimplicity/CloudOps/internal/k8s/service"
	promApi "github.com/GoSimplicity/CloudOps/internal/prometheus/api"
	promDao "github.com/GoSimplicity/CloudOps/internal/prometheus/dao"
	promService "github.com/GoSimplicity/CloudOps/internal/prometheus/service"
	"github.com/GoSimplicity/CloudOps/internal/startup"
	systemApi "github.com/GoSimplicity/CloudOps/internal/system/api"
	systemDao "github.com/GoSimplicity/CloudOps/internal/system/dao"
	systemService "github.com/GoSimplicity/CloudOps/internal/system/service"
	workorderApi "github.com/GoSimplicity/CloudOps/internal/workorder/api"
	workorderDao "github.com/GoSimplicity/CloudOps/internal/workorder/dao"
	workorderService "github.com/GoSimplicity/CloudOps/internal/workorder/service"
	ijwt "github.com/GoSimplicity/CloudOps/pkg/jwt"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var HandlerSet = wire.NewSet(
	systemApi.NewUserHandler,
	systemApi.NewRoleHandler,
	systemApi.NewMenuHandler,
	k8sApi.NewK8sClusterHandler,
	k8sApi.NewK8sNamespaceHandler,
	k8sApi.NewK8sDeploymentHandler,
	k8sApi.NewK8sPodHandler,
	k8sApi.NewK8sSvcHandler,
	k8sApi.NewK8sConfigMapHandler,
	promApi.NewScrapePoolHandler,
	promApi.NewAlertRuleHandler,
	promApi.NewAlertEventHandler,
	promApi.NewSendGroupHandler,
	workorderApi.NewTemplateHandler,
	workorderApi.NewInstanceHandler,
	workorderApi.NewFlowHandler,
	workorderApi.NewCommentHandler,
	aiApi.NewAIHandler,
)

var ServiceSet = wire.NewSet(
	systemService.NewUserService,
	systemService.NewRoleService,
	systemService.NewMenuService,
	k8sService.NewClusterService,
	k8sService.NewNamespaceService,
	k8sService.NewDeploymentService,
	k8sService.NewPodService,
	k8sService.NewServiceService,
	k8sService.NewConfigMapService,
	promService.NewScrapePoolService,
	promService.NewAlertRuleService,
	promService.NewAlertEventService,
	promService.NewSendGroupService,
	workorderService.NewTemplateService,
	workorderService.NewInstanceService,
	workorderService.NewFlowService,
	workorderService.NewCommentService,
	aiService.NewLLMService,
	aiService.NewRAGService,
)

var DaoSet = wire.NewSet(
	systemDao.NewUserDAO,
	systemDao.NewRoleDAO,
	systemDao.NewMenuDAO,
	k8sDao.NewClusterDAO,
	promDao.NewScrapePoolDAO,
	promDao.NewAlertRuleDAO,
	promDao.NewAlertEventDAO,
	promDao.NewSendGroupDAO,
	workorderDao.NewTemplateDAO,
	workorderDao.NewInstanceDAO,
	workorderDao.NewFlowDAO,
	workorderDao.NewCommentDAO,
	aiDao.NewKnowledgeDAO,
)

var InfraSet = wire.NewSet(
	InitLogger,
	InitDB,
	InitRedis,
	InitCasbin,
	InitMiddlewares,
	InitGinServer,
)

var ClientSet = wire.NewSet(
	k8sClient.NewK8sClient,
	vector.NewStore,
)

func InitApp() (*gin.Engine, *startup.ApplicationBootstrap, error) {
	wire.Build(
		InfraSet,
		HandlerSet,
		ServiceSet,
		DaoSet,
		ClientSet,
		ijwt.NewJWTHandler,
		startup.NewApplicationBootstrap,
	)
	return nil, nil, nil
}
