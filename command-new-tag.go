package main

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	. "gopkg.in/src-d/go-git.v4/_examples"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

type NewTagCommand struct {
}

func (c *NewTagCommand) Execute(_ []string) error {
	r, err := git.PlainOpen("")
	if err != nil {
		return err
	}

	hash, sm, err := findLastSemverRef(r)
	if err != nil {
		return err
	}
	logrus.Infof("Founded previous version `%s` in commit `%s`", sm.String(), hash)

	cIter, err := r.Log(&git.LogOptions{})
	if err != nil {
		return fmt.Errorf("can't get git log: %s", err)
	}

	// ... just iterates over the commits, printing it
	err = cIter.ForEach(func(c *object.Commit) error {
		var cc Change
		if err := cc.UnmarshalString(c.Message); err != nil {
			return err
		}
		fmt.Printf("%s\n", c.Hash)
		return nil
	})
	CheckIfError(err)

	return nil
}

func findLastSemverRef(r *git.Repository) (commitHash string, sm *semver.Version, err error) {
	var ref *plumbing.Reference
	tagrefs, err := r.Tags()
	if err != nil {
		return
	}
	err = tagrefs.ForEach(func(t *plumbing.Reference) error {
		s, err := semver.NewVersion(t.Name().Short())
		if err != nil {
			logrus.Warn(err)
		} else {
			if sm == nil || s.GreaterThan(sm) {
				sm = s
				ref = t
			}
		}
		return nil
	})
	if err != nil {
		return
	}

	rev, err := r.ResolveRevision(plumbing.Revision(ref.Name().String()))
	if err != nil {
		return
	}
	return rev.String(), sm, nil
}


/*
local prev_tag=$(git describe --abbrev=0 --tags)
    local messages=$(git log ${prev_tag}..HEAD --pretty=format:%s)

    if [[ -z "$messages" ]]; then
        bake_echo_red "Со времен ${prev_tag} новых коммитов не было"
        exit
    fi

    bake_echo_green "Предыдущий тег - ${prev_tag}"
    bake_echo_green "Коммиты:"
    echo -e "$messages"

    while true; do
        read -p $'  \e[36mВведите номер тега. Например, 0.0.92: \e[0m' new_tag;
        if [[ "$new_tag" =~ ^[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
            break
        else
            bake_echo_red "Слишком скучный номер версии"
        fi
    done

    local msg_one_line=$(echo "${messages}" | tr '\n' '|' | sed 's/|/ | /g')
    echo -e "  Создаем тег \e[32m${new_tag}\e[0m, с текстом \e[32m${msg_one_line}\e[0m"
    while true; do
        read -p "  [y/n]" -s -n 1 yn;
        case $yn in
            [Yy]* ) break;;
            [Nn]* ) exit;;
            * ) echo "  y или n";;
        esac
    done

    git tag -a -m "${messages}" ${new_tag} && bake_echo_green "Готово"
 */