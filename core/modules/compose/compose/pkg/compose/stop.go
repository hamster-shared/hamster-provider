/*
   Copyright 2020 Docker Compose CLI authors

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package compose

import (
	"context"
	"strings"

	"github.com/docker/compose/v2/pkg/api"
	"github.com/docker/compose/v2/pkg/progress"
)

func (s *composeService) Stop(ctx context.Context, projectName string, options api.StopOptions) error {
	return progress.Run(ctx, func(ctx context.Context) error {
		return s.stop(ctx, strings.ToLower(projectName), options)
	})
}

func (s *composeService) stop(ctx context.Context, projectName string, options api.StopOptions) error {
	w := progress.ContextWriter(ctx)

	containers, project, err := s.actualState(ctx, projectName, options.Services)
	if err != nil {
		return err
	}

	return InReverseDependencyOrder(ctx, project, func(c context.Context, service string) error {
		containersToStop := containers.filter(isService(service)).filter(isNotOneOff)
		return s.stopContainers(ctx, w, containersToStop, options.Timeout)
	})
}
