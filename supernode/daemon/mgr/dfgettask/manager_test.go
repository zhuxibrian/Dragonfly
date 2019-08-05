/*
 * Copyright The Dragonfly Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dfgettask

import (
	"context"
	"testing"

	"github.com/dragonflyoss/Dragonfly/apis/types"
	"github.com/dragonflyoss/Dragonfly/pkg/errortypes"

	"github.com/go-check/check"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

func init() {
	check.Suite(&DfgetTaskMgrTestSuite{})
}

type DfgetTaskMgrTestSuite struct {
	manager *Manager
}

// SetUpSuite does common setup in the beginning of each test.
func (s *DfgetTaskMgrTestSuite) SetUpSuite(c *check.C) {
	s.manager, _ = NewManager()
}

func (s *DfgetTaskMgrTestSuite) TestDfgetTaskMgr(c *check.C) {
	clientID := "foo"
	taskID := "00c4e7b174af7ed61c414b36ef82810ac0c98142c03e5748c00e1d1113f3c882"

	// Add
	dfgetTask := &types.DfGetTask{
		CID:       clientID,
		Path:      "/peer/file/taskFileName",
		PieceSize: 4 * 1024 * 1024,
		TaskID:    taskID,
		PeerID:    "foo-192.168.10.11-1553838710990554281",
	}

	err := s.manager.Add(context.Background(), dfgetTask)
	c.Check(err, check.IsNil)

	// Get
	dt, err := s.manager.Get(context.Background(), clientID, taskID)
	c.Check(err, check.IsNil)
	c.Check(dt, check.DeepEquals, &types.DfGetTask{
		CID:       clientID,
		Path:      "/peer/file/taskFileName",
		PieceSize: 4 * 1024 * 1024,
		TaskID:    taskID,
		Status:    types.DfGetTaskStatusWAITING,
		PeerID:    "foo-192.168.10.11-1553838710990554281",
	})

	// UpdateStatus
	err = s.manager.UpdateStatus(context.Background(), clientID, taskID, types.DfGetTaskStatusSUCCESS)
	c.Check(err, check.IsNil)

	dt, err = s.manager.Get(context.Background(), clientID, taskID)
	c.Check(err, check.IsNil)
	c.Check(dt, check.DeepEquals, &types.DfGetTask{
		CID:       clientID,
		Path:      "/peer/file/taskFileName",
		PieceSize: 4 * 1024 * 1024,
		TaskID:    taskID,
		Status:    types.DfGetTaskStatusSUCCESS,
		PeerID:    "foo-192.168.10.11-1553838710990554281",
	})

	// Delete
	err = s.manager.Delete(context.Background(), clientID, taskID)
	c.Check(err, check.IsNil)

	_, err = s.manager.Get(context.Background(), clientID, taskID)
	c.Check(errortypes.IsDataNotFound(err), check.Equals, true)
}
