package gclient

import (
	"context"
	"fmt"
	"gitlab.corp.qunar.com/qgrpc-go/config"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/common/param"
	"gitlab.corp.qunar.com/tcdev/qconfig-go/qconfig"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"testing"
)

const (
	testAppCode = "ops_watcher_gwhb"
	testToken   = "BvetjV3K+Mj6FJ6Qigy+yfc9AyTZfjA/0vlOEq1ZlvhF/csWuT77AN7PMtnT4H7IgiEdT0WlYDEhCTn922tY+HfuMSpzeOgMoTJbm0wpDpdKOgxN29AKf9vU69GjpLOEPTs7YHY1iC3DwuzEESCmrt7A0IW1/Eybxd4EstBFno4="
)

func TestServer(t *testing.T) {
	clientParam := param.QConfigClientParam{
		AppInfo: param.AppInfoParam{
			AppCode: testAppCode,
			Token:   testToken,
		},
	}

	err := qconfig.Init(clientParam, param.QConfigFiles{
		config.ClientConfigFile: &config.ClientConfig,
	})
	if err != nil {
		log.Fatalln(err)
	}

	client, err := New().
		WithContext(context.Background()).
		WithDialOptions(grpc.WithTransportCredentials(insecure.NewCredentials())).
		Build()

	fmt.Println(client.GetConn(testAppCode))
}
