// Copyright © 2018 Camunda Services GmbH (info@camunda.com)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package commands

import (
	"context"
	"github.com/camunda/camunda/clients/go/v8/internal/mock_pb"
	"github.com/camunda/camunda/clients/go/v8/internal/utils"
	"github.com/camunda/camunda/clients/go/v8/pkg/pb"
	"github.com/golang/mock/gomock"
	"testing"
)

func TestDeployResourceCommand_AddResourceFile(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	demoName := "../../cmd/zbctl/testdata/demo-process.bpmn"
	demoBytes := readBytes(t, demoName)
	anotherName := "../../cmd/zbctl/testdata/drg-force-user.dmn"
	anotherBytes := readBytes(t, anotherName)
	formName := "../../cmd/zbctl/testdata/deploy_form.form"
	formBytes := readBytes(t, formName)

	request := &pb.DeployResourceRequest{
		Resources: []*pb.Resource{
			{
				Name:    demoName,
				Content: demoBytes,
			},
			{
				Name:    anotherName,
				Content: anotherBytes,
			},
			{
				Name:    formName,
				Content: formBytes,
			},
		},
	}
	stub := &pb.DeployResourceResponse{}

	client.EXPECT().DeployResource(gomock.Any(), &utils.RPCTestMsg{Msg: request}).Return(stub, nil)

	command := NewDeployResourceCommand(client, func(context.Context, error) bool { return false })

	response, err := command.
		AddResourceFile(demoName).
		AddResourceFile(anotherName).
		AddResourceFile(formName).
		Send(context.Background())

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestDeployResourceCommand_AddResource(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	demoName := "../../cmd/zbctl/testdata/demo-process.bpmn"
	demoBytes := readBytes(t, demoName)

	request := &pb.DeployResourceRequest{
		Resources: []*pb.Resource{
			{
				Name:    demoName,
				Content: demoBytes,
			},
		},
	}
	stub := &pb.DeployResourceResponse{}

	client.EXPECT().DeployResource(gomock.Any(), &utils.RPCTestMsg{Msg: request}).Return(stub, nil)

	command := NewDeployResourceCommand(client, func(context.Context, error) bool { return false })

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultTestTimeout)
	defer cancel()

	response, err := command.
		AddResource(demoBytes, demoName).
		Send(ctx)

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}

func TestDeployResourceCommand_TenantId(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	client := mock_pb.NewMockGatewayClient(ctrl)

	demoName := "../../cmd/zbctl/testdata/demo-process.bpmn"
	demoBytes := readBytes(t, demoName)

	request := &pb.DeployResourceRequest{
		Resources: []*pb.Resource{
			{
				Name:    demoName,
				Content: demoBytes,
			},
		},
		TenantId: "1234",
	}
	stub := &pb.DeployResourceResponse{}

	client.EXPECT().DeployResource(gomock.Any(), &utils.RPCTestMsg{Msg: request}).Return(stub, nil)

	command := NewDeployResourceCommand(client, func(context.Context, error) bool { return false })

	ctx, cancel := context.WithTimeout(context.Background(), utils.DefaultTestTimeout)
	defer cancel()

	response, err := command.
		AddResource(demoBytes, demoName).
		TenantId("1234").
		Send(ctx)

	if err != nil {
		t.Errorf("Failed to send request")
	}

	if response != stub {
		t.Errorf("Failed to receive response")
	}
}
