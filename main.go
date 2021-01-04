package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
)

var err error

func main() {

	var restartingAuth = make(chan struct{}, 1)

	os.Remove("auth_service_log_file")
	// define the logger output file path
	serviceLogger, err := os.OpenFile("auth_service_log_file", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Error opening file: %v", err)
	}

	defer serviceLogger.Close()

	// set the log output
	log.SetOutput(serviceLogger)

	// setting the debug exec command
	log.Println("[GLOBAL][INFO] Creating the executable command")
	auth := exec.Command("go", "run", `C:\Users\Fakhri\go\src\github.com\fakhripraya\authentication-service\main.go`)
	whatsapp := exec.Command("go", "run", `C:\Users\Fakhri\go\src\github.com\fakhripraya\whatsapp-service\main.go`)
	email := exec.Command("go", "run", `C:\Users\Fakhri\go\src\github.com\fakhripraya\emailing-service\main.go`)

	// Whatsapp-service
	// Programmatically execute on develop whatsapp-service
	go func() {
		for {

			// set the stdout and stderr
			var out bytes.Buffer
			var stderr bytes.Buffer
			whatsapp.Stdout = &out
			whatsapp.Stderr = &stderr

			// Start the executable command
			log.Println("[WA][INFO] Execute the command to start the whatsapp server")
			err = whatsapp.Start()
			if err != nil {

				// If execute command fails, log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[WA][ERROR] Error occurred while starting the whatsapp server: ", err.Error())

				return
			}

			// Wait the program exit
			log.Println("[WA][INFO] The whatsapp server has been successfully started")
			err = whatsapp.Wait()
			if err != nil {

				// If exit code is not zero (0), log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[WA][ERROR] Error occurred while waiting the whatsapp server: ", err.Error())

				return
			}

			err = auth.Process.Signal(os.Kill)
			if err != nil {

				// If exit code is not zero (0), log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[WA][ERROR] Error occurred while sending signal to the authentication server: ", err.Error())

				return
			}

			// recieve the signal from the auth server
			log.Println("[WA][INFO] Waiting for the signal from authentication server")
			<-restartingAuth

			// If exit code is zero (0), restart the server
			log.Println("[WA][INFO] Restarting the whatsapp server")

		}
	}()

	// Emailing-service
	// Programmatically execute on develop emailing-service
	go func() {
		for {

			// set the stdout and stderr
			var out bytes.Buffer
			var stderr bytes.Buffer
			email.Stdout = &out
			email.Stderr = &stderr

			// Start the executable command
			log.Println("[EMAIL][INFO] Execute the command to start the email server")
			err = email.Start()
			if err != nil {

				// If execute command fails, log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[EMAIL][ERROR] Error occurred while starting the email server: ", err.Error())

				return
			}

			// Wait the program exit
			log.Println("[EMAIL][INFO] The email server has been successfully started")
			err = email.Wait()
			if err != nil {

				// If exit code is not zero (0), log the errors
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				fmt.Println("Result: " + out.String())
				log.Fatalln("[EMAIL][ERROR] Error occurred while waiting the email server: ", err.Error())

				return
			}

			err = auth.Process.Signal(os.Kill)
			if err != nil {

				// If exit code is not zero (0), log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[EMAIL][ERROR] Error occurred while sending signal to the authentication server: ", err.Error())

				return
			}

			// recieve the signal from the auth server
			log.Println("[EMAIL][INFO] Waiting for the signal from authentication server")
			<-restartingAuth

			// If exit code is zero (0), restart the server
			log.Println("[EMAIL][INFO] Restarting the email server")

		}
	}()

	// Authentication-service
	// Programmatically execute on develop Authentication-service
	go func() {
		for {

			// set the stdout and stderr
			var out bytes.Buffer
			var stderr bytes.Buffer
			auth.Stdout = &out
			auth.Stderr = &stderr

			// Start the executable command
			log.Println("[AUTH][INFO] Execute the command to start the authentication server")
			err = auth.Start()
			if err != nil {

				// If execute command fails, log the errors
				log.Println(fmt.Sprint(err) + ": " + stderr.String())
				log.Println("Result: " + out.String())
				log.Fatalln("[AUTH][ERROR] Error occurred while starting the authentication server: ", err.Error())

				return
			}

			// signaling the whatsapp server
			log.Println("[AUTH][INFO] Signaling the whatsapp server")
			restartingAuth <- struct{}{}

			// Wait the program exit
			log.Println("[AUTH][INFO] The authentication server has been successfully started")
			err = auth.Wait()
			if err != nil {

				// If exit code is not zero (0), log the errors
				fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
				fmt.Println("Result: " + out.String())
				log.Fatalln("[AUTH][ERROR] Error occurred while starting the authentication server: ", err.Error())

				return
			}

			// If exit code is zero (0), restart the server
			log.Println("[AUTH][INFO] Restarting the authentication server")

		}
	}()

	// trap sigterm or interrupt and gracefully shutdown the server
	channel := make(chan os.Signal, 1)
	signal.Notify(channel, os.Interrupt)
	signal.Notify(channel, os.Kill)

	// Block until a signal is received.
	sig := <-channel
	log.Println("Got signal : ", sig)

	os.Exit(0)
}
