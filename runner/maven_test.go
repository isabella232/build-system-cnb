/*
 * Copyright 2018-2020 the original author or authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      https://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package runner_test

import (
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/build-system-cnb/buildsystem"
	"github.com/cloudfoundry/build-system-cnb/runner"
	"github.com/cloudfoundry/libcfbuildpack/buildpackplan"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/onsi/gomega"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestMaven(t *testing.T) {
	spec.Run(t, "Maven", func(t *testing.T, when spec.G, it spec.S) {

		g := gomega.NewWithT(t)

		var f *test.BuildFactory

		it.Before(func() {
			f = test.NewBuildFactory(t)

			f.AddDependency(buildsystem.MavenDependency, filepath.Join("testdata", "stub-maven.tar.gz"))
			f.AddPlan(buildpackplan.Plan{Name: buildsystem.MavenDependency})
			test.TouchFile(t, f.Build.Application.Root, ".mvn")
			test.TouchFile(t, f.Build.Application.Root, "mvnw")
		})

		when("working with JAR file", func() {

			it.Before(func() {
				test.CopyFile(t, filepath.Join("testdata", "stub-executable.jar"),
					filepath.Join(f.Build.Application.Root, "target", "stub-executable.jar"))
			})

			it("builds application", func() {
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				g.Expect(f.Runner.Commands[1]).
					To(gomega.Equal(test.Command{
						Bin:  filepath.Join(f.Build.Application.Root, "mvnw"),
						Dir:  f.Build.Application.Root,
						Args: []string{"-Dmaven.test.skip=true", "package"},
					}))
			})

			it("builds application with custom command", func() {
				defer test.ReplaceEnv(t, "BP_BUILD_ARGUMENTS", "test configured arguments")()
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				g.Expect(f.Runner.Commands[1]).
					To(gomega.Equal(test.Command{
						Bin:  filepath.Join(f.Build.Application.Root, "mvnw"),
						Dir:  f.Build.Application.Root,
						Args: []string{"test", "configured", "arguments"},
					}))
			})

			it("removes source code", func() {
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				g.Expect(f.Build.Application.Root).To(gomega.BeADirectory())
				g.Expect(filepath.Join(f.Build.Application.Root, ".mvn")).NotTo(gomega.BeAnExistingFile())
				g.Expect(filepath.Join(f.Build.Application.Root, "mvnw")).NotTo(gomega.BeAnExistingFile())
				g.Expect(filepath.Join(f.Build.Application.Root, "target")).NotTo(gomega.BeAnExistingFile())
			})

			it("explodes built application", func() {
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("build-system-application")
				g.Expect(layer).To(test.HaveLayerMetadata(false, false, false))
				g.Expect(filepath.Join(f.Build.Application.Root, "fixture-marker")).To(gomega.BeARegularFile())
			})
		})

		when("working with WAR file", func() {

			it.Before(func() {
				test.CopyFile(t, filepath.Join("testdata", "stub-application.war"),
					filepath.Join(f.Build.Application.Root, "target", "stub-application.war"))
			})

			it("explodes built application", func() {
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("build-system-application")
				g.Expect(layer).To(test.HaveLayerMetadata(false, false, false))
				g.Expect(filepath.Join(f.Build.Application.Root, "fixture-marker")).To(gomega.BeARegularFile())
			})
		})

		when("working with modules", func() {

			it.Before(func() {
				test.CopyFile(t, filepath.Join("testdata", "stub-executable.jar"),
					filepath.Join(f.Build.Application.Root, "test-module", "target", "stub-executable.jar"))
			})

			it("explodes built application", func() {
				defer test.ReplaceEnv(t, "BP_BUILT_MODULE", "test-module")()
				f.Runner.Outputs = []string{"test-java-version"}

				b, _, err := buildsystem.NewMavenBuildSystem(f.Build)
				g.Expect(err).NotTo(gomega.HaveOccurred())
				r, err := runner.NewMavenRunner(f.Build, b)
				g.Expect(err).NotTo(gomega.HaveOccurred())

				g.Expect(r.Contribute()).To(gomega.Succeed())

				layer := f.Build.Layers.Layer("build-system-application")
				g.Expect(layer).To(test.HaveLayerMetadata(false, false, false))
				g.Expect(filepath.Join(f.Build.Application.Root, "fixture-marker")).To(gomega.BeARegularFile())
			})
		})
	}, spec.Report(report.Terminal{}))
}
