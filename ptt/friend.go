package ptt

import (
	"bufio"
	"fmt"
	"io"
	"math/rand"
	"os"

	"github.com/Ptt-official-app/go-pttbbs/cmbbs/names"
	"github.com/Ptt-official-app/go-pttbbs/cmbbs/path"

	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/Ptt-official-app/go-pttbbs/types"

	log "github.com/sirupsen/logrus"
)

func friendDeleteAll(userID *ptttype.UserID_t, friendType ptttype.FriendType) error {
	filename, err := path.SetHomeFile(userID, ptttype.FriendFile[friendType])
	if err != nil { //unable to get the file. assuming not-exists
		return err
	}

	file, err := os.Open(filename)
	if err != nil { //unable to open the file. assuming not-exists
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	for line, err := types.ReadLine(reader); err == nil; line, err = types.ReadLine(reader) {
		friendID := &ptttype.UserID_t{}
		copy(friendID[:], line)
		if !names.IsValidUserID(friendID) {
			continue
		}

		// XXX race-condition.
		deleteUserFriend(friendID, userID, friendType) // remove me from my friend.
	}
	return nil
}

func deleteUserFriend(userID *ptttype.UserID_t, friendID *ptttype.UserID_t, friendType ptttype.FriendType) {
	if types.Cstrcmp(userID[:], friendID[:]) == 0 {
		return
	}

	//XXX If I remove the user-id from rejected,
	//    I still remove they from ALOHA list.
	//
	//    deleteUserFriend used only in friendDeleteAll,
	//    friendDeleteAll used only in killUser, with type FRIEND_ALOHA
	filename, err := path.SetHomeFile(userID, ptttype.FN_ALOHA)
	if err != nil {
		return
	}

	_ = deleteFriendFromFile(filename, friendID, false)
}

func deleteFriendFromFile(filename string, friend *ptttype.UserID_t, isCaseSensitive bool) bool {
	// XXX race-condition
	randNum := rand.Intn(0xfff)
	randStr := fmt.Sprintf("%3.3X", randNum)
	new_filename := filename + "." + randStr

	file, err := os.OpenFile(filename, os.O_RDONLY, 0600)
	if err != nil { // no file
		return false
	}
	defer file.Close()

	new_file, err := os.OpenFile(new_filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil { // unable to create new file.
		return false
	}
	defer new_file.Close()

	userIDInFile := &ptttype.UserID_t{}
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(new_file)

	var line []byte
	for {
		line, err = types.ReadLine(reader)
		if err != nil {
			if err == io.EOF {
				err = nil // make it clear about err here
				break
			}
			log.Errorf("friend.deleteFriendFromFile: unable to read file: filename: %v new_filename: %v e: %v", filename, new_filename, err)
			return false
		}

		copy(userIDInFile[:], line[:])
		if isCaseSensitive && (types.Cstrcmp(friend[:], userIDInFile[:]) == 0) ||
			!isCaseSensitive && (types.Cstrcasecmp(friend[:], userIDInFile[:]) == 0) {
			sanitizedUserIDInFile := types.CstrToBytes(userIDInFile[:])
			_, err = writer.Write(sanitizedUserIDInFile)
			if err != nil { // unable to write-bytes into tmp-file.
				log.Errorf("friend.deleteFriendFromFile: unable write to tmp-file (possible zombie-tmp-file): filename: %v new_filename: %v e: %v", filename, new_filename, err)
				return false
			}
			err = writer.WriteByte('\n')
			if err != nil { // unable to write new-line into tmp-file.
				log.Errorf("friend.deleteFriendFromFile: unable write to tmp-file (possible zombie-tmp-file): filename: %v new_filename: %v e: %v", filename, new_filename, err)
				return false
			}

		}
	}

	err = os.Rename(new_filename, filename)
	if err != nil {
		return false
	}

	return true
}
