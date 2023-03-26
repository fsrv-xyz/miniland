package account

import (
	"errors"
	"os"
	"path"
	"sync"
)

type State struct {
	passwdFile *passwdFile
	groupFile  *groupFile
	mux        sync.Mutex
}

const (
	PasswdFile = "/etc/passwd"
	GroupFile  = "/etc/group"
)

func NewState() *State {
	return &State{
		passwdFile: &passwdFile{},
		groupFile:  &groupFile{},
	}
}
func (state *State) AddGroup(options ...GroupOption) error {
	// create a new group with default values
	group := &Group{
		groupEntry: groupEntry{
			Password: "x",
		},
	}
	for _, option := range options {
		if optionApplyError := option(group); optionApplyError != nil {
			return optionApplyError
		}
	}

	// validate the group
	if group.Name == "" {
		return errors.New("group name is required")
	}

	state.mux.Lock()
	defer state.mux.Unlock()
	if state.groupFile == nil {
		state.groupFile = &groupFile{}
	}
	state.groupFile.Entries = append(state.groupFile.Entries, group.groupEntry)
	return nil
}

func (state *State) AddUser(options ...UserOption) error {
	// create a new user with default values
	user := &User{
		passwdEntry: passwdEntry{
			Password: "x",
			Shell:    "/sbin/nologin",
		},
	}
	for _, option := range options {
		if optionApplyError := option(user); optionApplyError != nil {
			return optionApplyError
		}
	}

	// validate the user
	if user.Username == "" {
		return errors.New("user name is required")
	}

	// create the group if requested
	if user.createGroup {
		if groupCreateError := state.AddGroup(user.groupOptions...); groupCreateError != nil {
			return groupCreateError
		}
	}

	state.mux.Lock()
	defer state.mux.Unlock()
	if state.passwdFile == nil {
		state.passwdFile = &passwdFile{}
	}
	state.passwdFile.Entries = append(state.passwdFile.Entries, user.passwdEntry)
	return nil
}

func (state *State) WriteFiles() error {
	state.mux.Lock()
	defer state.mux.Unlock()

	var result error

	for _, cfgFile := range []string{PasswdFile, GroupFile} {
		resultingError := func() error {
			tmpFile, tmpFileCreateError := os.CreateTemp(path.Dir(cfgFile), path.Base(cfgFile)+".*")
			if tmpFileCreateError != nil {
				return tmpFileCreateError
			}

			var encodeError error
			switch cfgFile {
			case PasswdFile:
				encodeError = encoder(tmpFile, state.passwdFile.Entries)
			case GroupFile:
				encodeError = encoder(tmpFile, state.groupFile.Entries)
			}
			if encodeError != nil {
				return encodeError
			}

			if syncError := tmpFile.Sync(); syncError != nil {
				return syncError
			}

			if closeError := tmpFile.Close(); closeError != nil {
				return closeError
			}

			return os.Rename(tmpFile.Name(), cfgFile)
		}()

		if resultingError != nil {
			result = errors.Join(result, resultingError)
		}
	}
	return result
}

func (state *State) ReadFiles() error {
	state.mux.Lock()
	defer state.mux.Unlock()

	var result error

	for _, cfgFile := range []string{PasswdFile, GroupFile} {
		resultingError := func() error {
			file, fileOpenError := os.Open(cfgFile)
			if fileOpenError != nil {
				return fileOpenError
			}
			defer file.Close()

			var decodeError error
			switch cfgFile {
			case PasswdFile:
				var content []passwdEntry
				if content, decodeError = decoder[passwdEntry](file, passwdEntry{}.parse); decodeError == nil {
					state.passwdFile.Entries = content
				}
			case GroupFile:
				var content []groupEntry
				if content, decodeError = decoder[groupEntry](file, groupEntry{}.parse); decodeError == nil {
					state.groupFile.Entries = content
				}
			}

			return decodeError
		}()

		if resultingError != nil {
			result = errors.Join(result, resultingError)
		}
	}
	return result
}
