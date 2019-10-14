package main

import (
	"net/http"
	"reflect"
	"testing"

	"cicd-github.quickplay.com/platform-engineering/metadata-mgmt-services/pkg/repository"
	"github.com/gorilla/mux"
)

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func Test_initService(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			initService()
		})
	}
}

func Test_initDatabase(t *testing.T) {
	tests := []struct {
		name string
		want repository.DbRepository
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initDatabase(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initDatabase() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_initHttpServer(t *testing.T) {
	type args struct {
		r *mux.Router
	}
	tests := []struct {
		name string
		args args
		want *http.Server
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := initHttpServer(tt.args.r); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("initHttpServer() = %v, want %v", got, tt.want)
			}
		})
	}
}
