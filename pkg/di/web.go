package di

import (
	_ "github.com/GoSimplicity/CloudOps/docs"
	aiApi "github.com/GoSimplicity/CloudOps/internal/ai/api"
	k8sApi "github.com/GoSimplicity/CloudOps/internal/k8s/api"
	promApi "github.com/GoSimplicity/CloudOps/internal/prometheus/api"
	systemApi "github.com/GoSimplicity/CloudOps/internal/system/api"
	workorderApi "github.com/GoSimplicity/CloudOps/internal/workorder/api"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func InitGinServer(
	m []gin.HandlerFunc,
	// system
	userHdl *systemApi.UserHandler,
	roleHdl *systemApi.RoleHandler,
	menuHdl *systemApi.MenuHandler,
	// k8s
	clusterHdl *k8sApi.K8sClusterHandler,
	nsHdl *k8sApi.K8sNamespaceHandler,
	deployHdl *k8sApi.K8sDeploymentHandler,
	podHdl *k8sApi.K8sPodHandler,
	svcHdl *k8sApi.K8sSvcHandler,
	cmHdl *k8sApi.K8sConfigMapHandler,
	// prometheus
	scrapePoolHdl *promApi.ScrapePoolHandler,
	alertRuleHdl *promApi.AlertRuleHandler,
	alertEventHdl *promApi.AlertEventHandler,
	sendGroupHdl *promApi.SendGroupHandler,
	// workorder
	templateHdl *workorderApi.TemplateHandler,
	instanceHdl *workorderApi.InstanceHandler,
	flowHdl *workorderApi.FlowHandler,
	commentHdl *workorderApi.CommentHandler,
	// ai
	aiHdl *aiApi.AIHandler,
) *gin.Engine {
	server := gin.Default()
	server.Use(m...)

	server.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	server.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })

	userHdl.RegisterRouters(server)
	roleHdl.RegisterRouters(server)
	menuHdl.RegisterRouters(server)
	clusterHdl.RegisterRouters(server)
	nsHdl.RegisterRouters(server)
	deployHdl.RegisterRouters(server)
	podHdl.RegisterRouters(server)
	svcHdl.RegisterRouters(server)
	cmHdl.RegisterRouters(server)
	scrapePoolHdl.RegisterRouters(server)
	alertRuleHdl.RegisterRouters(server)
	alertEventHdl.RegisterRouters(server)
	sendGroupHdl.RegisterRouters(server)
	templateHdl.RegisterRouters(server)
	instanceHdl.RegisterRouters(server)
	flowHdl.RegisterRouters(server)
	commentHdl.RegisterRouters(server)
	aiHdl.RegisterRouters(server)

	return server
}
