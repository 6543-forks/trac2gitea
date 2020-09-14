// Copyright 2020 Steve Jefferson. All rights reserved.
// Use of this source code is governed by a GPL-style
// license that can be found in the LICENSE file.

package importer

import (
	"github.com/stevejefferson/trac2gitea/accessor/gitea"
	"github.com/stevejefferson/trac2gitea/accessor/trac"
)

func truncateString(str string, maxlen int) string {
	strLen := len(str)
	if strLen > maxlen {
		return str[0:maxlen] + "..."
	}
	return str
}

// importTicketComment imports a single ticket comment from Trac to Gitea, returns ID of created comment or -1 if comment already exists
func (importer *Importer) importTicketComment(issueID int64, tracComment *trac.TicketComment, commentTime int64, userMap map[string]string) (int64, error) {
	authorID, err := importer.getUser(tracComment.Author, userMap)
	if err != nil {
		return -1, err
	}

	// record Trac comment author as original author if it cannot be mapped onto a Gitea user
	originalAuthorName := ""
	if authorID == -1 {
		authorID = importer.defaultAuthorID
		originalAuthorName = tracComment.Author
	}

	convertedText := importer.markdownConverter.TicketConvert(tracComment.TicketID, tracComment.Text)
	giteaComment := gitea.IssueComment{AuthorID: authorID, OriginalAuthorID: 0, OriginalAuthorName: originalAuthorName, Text: convertedText, Time: commentTime}
	commentID, err := importer.giteaAccessor.AddIssueComment(issueID, &giteaComment)
	if err != nil {
		return -1, err
	}

	// add association between issue and comment author
	err = importer.giteaAccessor.AddIssueUser(issueID, authorID)
	if err != nil {
		return -1, err
	}

	return commentID, nil
}

func (importer *Importer) importTicketChanges(ticketID int64, issueID int64, lastUpdate int64, userMap map[string]string) (int64, error) {
	commentLastUpdate := lastUpdate
	err := importer.tracAccessor.GetTicketChanges(ticketID, func(change *trac.TicketChange) error {
		switch change.ChangeType {
		case trac.TicketCommentType:
			commentID, err := importer.importTicketComment(issueID, change.Comment, change.Time, userMap)
			if err != nil {
				return err
			}

			if commentID != -1 && commentLastUpdate < change.Time {
				commentLastUpdate = change.Time
			}
		}

		return nil
	})
	if err != nil {
		return 0, err
	}

	return commentLastUpdate, nil
}