/*
 * Copyright 2018 the original author or authors.
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

package runner_test

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/build-system-buildpack/buildsystem"
	"github.com/cloudfoundry/build-system-buildpack/runner"
	"github.com/cloudfoundry/libcfbuildpack/layers"
	"github.com/cloudfoundry/libcfbuildpack/test"
	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"
)

func TestGradle(t *testing.T) {
	spec.Run(t, "Gradle", testGradle, spec.Report(report.Terminal{}))
}

func testGradle(t *testing.T, when spec.G, it spec.S) {

	it("builds application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, buildsystem.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, buildsystem.GradleDependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "gradlew")

		g, _, err := buildsystem.NewGradleBuildSystem(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := runner.NewGradleRunner(f.Build, g)

		e := &testExecutor{Outputs: []string{"test-java-version"}}
		r.Executor = e

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		expected := []string{filepath.Join(f.Build.Application.Root, "gradlew"), "-x", "test", "build"}
		if !reflect.DeepEqual(e.Commands[1].Args, expected) {
			t.Errorf("Cmd.Args = %s, expected %s", e.Commands[1].Args, expected)
		}
	})

	it("removes source code", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, buildsystem.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, buildsystem.GradleDependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "gradlew")

		g, _, err := buildsystem.NewGradleBuildSystem(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := runner.NewGradleRunner(f.Build, g)

		e := &testExecutor{Outputs: []string{"test-java-version"}}
		r.Executor = e

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		fi, err := os.Lstat(f.Build.Application.Root)
		if err != nil {
			t.Fatal(err)
		}

		if fi.Mode()&os.ModeSymlink != os.ModeSymlink {
			t.Errorf("Application.Root.Mode() = %s, expected symlink", fi.Mode())
		}
	})

	it("explodes built application", func() {
		f := test.NewBuildFactory(t)
		f.AddDependency(t, buildsystem.GradleDependency, "stub-gradle.zip")
		f.AddBuildPlan(t, buildsystem.GradleDependency, buildplan.Dependency{})

		test.TouchFile(t, f.Build.Application.Root, "gradlew")

		g, _, err := buildsystem.NewGradleBuildSystem(f.Build)
		if err != nil {
			t.Fatal(err)
		}

		r := runner.NewGradleRunner(f.Build, g)

		e := &testExecutor{Outputs: []string{"test-java-version"}}
		r.Executor = e

		source := test.FixturePath(t, "stub-application.jar")
		destination := filepath.Join(f.Build.Application.Root, "build", "libs", "stub-application.jar")
		if err := layers.CopyFile(source, destination); err != nil {
			t.Fatal(err)
		}

		if err := r.Contribute(); err != nil {
			t.Fatal(err)
		}

		layer := f.Build.Layers.Layer("build-system-application")
		test.BeLayerLike(t, layer, true, false, true)
		exists, err := layers.FileExists(filepath.Join(layer.Root, "fixture-marker"))
		if err != nil {
			t.Fatal(err)
		}

		if !exists {
			t.Errorf("Expected application to be expanded, but was not")
		}
	})
}
