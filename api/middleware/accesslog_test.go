package middleware

import (
	"net/http"
	"reflect"
	"testing"
)

/*
 * Auto generated by Quickstart. 
 * Developer to replace the generated test cases with real ones.
 */
func TestAccessLog(t *testing.T) {
	type args struct {
		next http.Handler
	}
	tests := []struct {
		name string
		args args
		want http.Handler
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AccessLog(tt.args.next); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AccessLog() = %v, want %v", got, tt.want)
			}
		})
	}
}
