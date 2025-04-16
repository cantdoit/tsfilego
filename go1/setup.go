package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func main() {
	// Dynamically resolve the script location (directory of the running script file)
	scriptFile, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get the current working directory:", err)
		os.Exit(1)
	}

	// Calculate the project root (navigate up to the parent of `go1`)
	projectPath := filepath.Clean(filepath.Join(scriptFile, "..")) // Parent directory of `go1`

	// Paths for the C++ and Go project
	cppPath := filepath.Join(projectPath, "cpp")         // Absolute path to the C++ source folder
	buildPath := filepath.Join(cppPath, "build")         // Path for the CMake build directory
	includePath := filepath.Join(projectPath, "include") // Directory for include headers
	libOutputPath := filepath.Join(buildPath, "lib")     // Path for the built C++ libraries

	// Step 1: Ensure CMake Build Directory Exists
	if _, err := os.Stat(buildPath); os.IsNotExist(err) {
		err := os.MkdirAll(buildPath, os.ModePerm)
		if err != nil {
			fmt.Println("Failed to create build directory:", err)
			os.Exit(1)
		}
	}

	// Step 2: Run CMake to Configure the Build
	fmt.Println("Configuring CMake project...")
	cmd := exec.Command("cmake", cppPath, "-DCMAKE_BUILD_TYPE=Release") // Use absolute cppPath
	cmd.Dir = buildPath                                                 // Set build directory as the working dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("CMake configuration failed:", err)
		os.Exit(1)
	}

	// Step 3: Build the C++ Libraries with CMake
	fmt.Println("Building the C++ project...")
	cmd = exec.Command("cmake", "--build", ".", "--config", "Release") // Trigger actual build
	cmd.Dir = buildPath                                                // Use buildPath as the working directory
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		fmt.Println("CMake build failed:", err)
		os.Exit(1)
	}

	// Step 4: Ensure Headers Are Available
	headerFile := filepath.Join(cppPath, "src", "cwrapper", "TsFile-cwrapper.h") // Full path to header
	targetHeader := filepath.Join(includePath, "TsFile-cwrapper.h")              // Destination
	if _, err := os.Stat(headerFile); err == nil {
		fmt.Println("Copying header file to include directory...")
		if _, err := os.Stat(includePath); os.IsNotExist(err) {
			err = os.MkdirAll(includePath, os.ModePerm)
			if err != nil {
				fmt.Println("Failed to create include directory:", err)
				os.Exit(1)
			}
		}
		copyFile(headerFile, targetHeader)
	} else {
		fmt.Println("Header file not found:", headerFile)
		os.Exit(1)
	}

	// Print CGO Configuration
	fmt.Printf(`
C++ library build complete!

To use in Go:
    export CGO_CFLAGS=-I%s
    export CGO_LDFLAGS=-L%s -ltsfile
`, includePath, libOutputPath)
}

// Utility to copy files
func copyFile(source, target string) {
	input, err := os.ReadFile(source)
	if err != nil {
		fmt.Println("Failed to read file:", err)
		os.Exit(1)
	}
	err = os.WriteFile(target, input, 0644)
	if err != nil {
		fmt.Println("Failed to write file:", err)
		os.Exit(1)
	}
}
