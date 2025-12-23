// © 2025 Sharon Aicler (saichler@gmail.com)
//
// Layer 8 Ecosystem is licensed under the Apache License, Version 2.0.
// You may obtain a copy of the License at:
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package utils provides test helper functions for creating test data.
// It wraps the l8test infrastructure to create test model instances.

package utils

import (
	"github.com/saichler/l8test/go/infra/t_resources"
	"github.com/saichler/l8types/go/testtypes"
)

// CreateTestModelInstance creates a test protobuf instance with the given index.
// Delegates to the l8test infrastructure for consistent test data generation.
func CreateTestModelInstance(index int) *testtypes.TestProto {
	return t_resources.CreateTestModelInstance(index)
}
