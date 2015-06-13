package service

import (
	"bytes"
	"errors"
	"log"
	"net/http"
	"os/exec"

	r "github.com/dancannon/gorethink"
	"github.com/robertkrimen/otto"
)

type Service struct {
	Session *r.Session
}

type ExecuteArgs struct {
	Token  string
	Name   string
	Branch string
	Args   interface{}
}

type ExecuteReply struct {
	Response string
}

func (s *Service) Execute(req *http.Request, args *ExecuteArgs, reply *ExecuteReply) error {
	cursor, err := r.Table("tokens").Get(args.Token).Run(s.Session)
	if err != nil {
		return err
	}
	defer cursor.Close()
	var token *Token
	if err := cursor.One(&token); err != nil {
		return err
	}

	if token.Restriction != nil && len(token.Restriction) > 0 {
		found := false
		for _, s2 := range token.Restriction {
			if args.Name == s2 {
				found = true
				break
			}
		}

		if !found {
			return errors.New("Token does not allow accessing this script")
		}
	}

	cursor, err = r.Table("scripts").Get(args.Name).Run(s.Session)
	if err != nil {
		return err
	}
	defer cursor.Close()
	var script *Script
	if err := cursor.One(&script); err != nil {
		return err
	}

	vm := otto.New()
	vm.Set("exec", func(call otto.FunctionCall) otto.Value {
		sa := []string{}

		for _, arg := range call.ArgumentList {
			ss, err := arg.ToString()
			if err != nil {
				val, _ := vm.ToValue(map[string]interface{}{
					"error": err.Error(),
				})
				return val
			}

			sa = append(sa, ss)
		}

		if len(sa) > 0 {
			path, err := exec.LookPath(sa[0])
			if err != nil {
				val, _ := vm.ToValue(map[string]interface{}{
					"error": err.Error(),
				})
				return val
			}

			var cmd *exec.Cmd
			if len(sa) > 1 {
				cmd = exec.Command(path, sa[1:]...)
			} else {
				cmd = exec.Command(path)
			}

			var (
				oo = &bytes.Buffer{}
				oe = &bytes.Buffer{}
			)

			cmd.Stdout = oo
			cmd.Stderr = oe

			var val otto.Value

			err = cmd.Run()
			if err != nil {
				val, _ = vm.ToValue(map[string]interface{}{
					"error":  err.Error(),
					"stdout": oo.String(),
					"stderr": oe.String(),
				})
			} else {
				val, _ = vm.ToValue(map[string]interface{}{
					"error":  nil,
					"stdout": oo.String(),
					"stderr": oe.String(),
				})
			}

			log.Print(oo.String())
			log.Print(oe.String())

			return val
		}

		val, _ := vm.ToValue(map[string]interface{}{
			"error": "Invalid arguments",
		})
		return val
	})

	vm.Set("args", args.Args)

	value, err := vm.Run(script.Code)
	if err != nil {
		return err
	}
	response, err := value.ToString()
	if err != nil {
		return err
	}
	reply.Response = response

	return nil
}
