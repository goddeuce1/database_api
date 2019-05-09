package database

import (
	"fmt"
	"io/ioutil"

	middlewares "park_base/park_db/sqlops"

	"github.com/jackc/pgx"
)

//Application - application struct
type Application struct {
	DB *pgx.ConnPool
}

//App - export
var App Application

//OpenConnection - connects to database
func (a *Application) OpenConnection(input string) {
	pgxConfig, _ := pgx.ParseURI(input)

	a.DB, _ = pgx.NewConnPool(
		pgx.ConnPoolConfig{
			ConnConfig:     pgxConfig,
			MaxConnections: 10,
		})

	if query, err := ioutil.ReadFile("database/sql/items.sql"); err != nil {
		fmt.Println(err)
	} else {
		if _, err := a.DB.Exec(string(query)); err != nil {
			fmt.Println(err)
		}
	}

	//USER
	a.DB.Prepare("UCMGetByNick", middlewares.UCMGetByNick)
	a.DB.Prepare("UCMGetByNickOrMail", middlewares.UCMGetByNickOrMail)
	a.DB.Prepare("UCMInsertValues", middlewares.UCMInsertValues)
	a.DB.Prepare("UPPUpdateSettings", middlewares.UPPUpdateSettings)

	//POSTS
	a.DB.Prepare("PIDUUpdateMessage", middlewares.PIDUUpdateMessage)
	a.DB.Prepare("PIDUGetPostByID", middlewares.PIDUGetPostByID)
	a.DB.Prepare("PIDUGetForumByName", middlewares.PIDUGetForumByName)
	a.DB.Prepare("PIDUGetThreadByID", middlewares.PIDUGetThreadByID)

	//THREADS
	a.DB.Prepare("TPSinceDescLimitTree", middlewares.TPSinceDescLimitTree)
	a.DB.Prepare("TPSinceDescLimitParentTree", middlewares.TPSinceDescLimitParentTree)
	a.DB.Prepare("TPSinceDescLimitFlat", middlewares.TPSinceDescLimitFlat)
	a.DB.Prepare("TPSinceAscLimitTree", middlewares.TPSinceAscLimitTree)
	a.DB.Prepare("TPSinceAscLimitParentTree", middlewares.TPSinceAscLimitParentTree)
	a.DB.Prepare("TPSinceAscLimitFlat", middlewares.TPSinceAscLimitFlat)
	a.DB.Prepare("TPDescLimitTree", middlewares.TPDescLimitTree)
	a.DB.Prepare("TPDescLimitParentTree", middlewares.TPDescLimitParentTree)
	a.DB.Prepare("TPDescLimitFlat", middlewares.TPDescLimitFlat)
	a.DB.Prepare("TPAscLimitTree", middlewares.TPAscLimitTree)
	a.DB.Prepare("TPAscLimitParentTree", middlewares.TPAscLimitParentTree)
	a.DB.Prepare("TPAscLimitFlat", middlewares.TPAscLimitFlat)

	a.DB.Prepare("TFByID", middlewares.TFByID)
	a.DB.Prepare("TFBySlug", middlewares.TFBySlug)
	a.DB.Prepare("TDPUpdateMessageID", middlewares.TDPUpdateMessageID)
	a.DB.Prepare("TSVoteByID", middlewares.TSVoteByID)
	a.DB.Prepare("TCMFindPostByParent", middlewares.TCMFindPostByParent)
	a.DB.Prepare("TCMUpdateForumPostsCount", middlewares.TCMUpdateForumPostsCount)

	//FORUMS
	a.DB.Prepare("FSTSelectThreadsLSD", middlewares.FSTSelectThreadsLSD)
	a.DB.Prepare("FSTSelectThreadsLD", middlewares.FSTSelectThreadsLD)
	a.DB.Prepare("FSTSelectThreadsLS", middlewares.FSTSelectThreadsLS)
	a.DB.Prepare("FSTSelectThreadsL", middlewares.FSTSelectThreadsL)
	a.DB.Prepare("FSTSelectThreadsSD", middlewares.FSTSelectThreadsSD)
	a.DB.Prepare("FSTSelectThreadsD", middlewares.FSTSelectThreadsD)
	a.DB.Prepare("FSTSelectThreadsS", middlewares.FSTSelectThreadsS)
	a.DB.Prepare("FSTSelectThreads", middlewares.FSTSelectThreads)

	a.DB.Prepare("FSTSelectUsersLSD", middlewares.FSTSelectUsersLSD)
	a.DB.Prepare("FSTSelectUsersLD", middlewares.FSTSelectUsersLD)
	a.DB.Prepare("FSTSelectUsersLS", middlewares.FSTSelectUsersLS)
	a.DB.Prepare("FSTSelectUsersL", middlewares.FSTSelectUsersL)
	a.DB.Prepare("FSTSelectUsersSD", middlewares.FSTSelectUsersSD)
	a.DB.Prepare("FSTSelectUsersD", middlewares.FSTSelectUsersD)
	a.DB.Prepare("FSTSelectUsersS", middlewares.FSTSelectUsersS)
	a.DB.Prepare("FSTSelectUsers", middlewares.FSTSelectUsers)

	a.DB.Prepare("FCMSelectNick", middlewares.FCMSelectNick)
	a.DB.Prepare("FCMInsertValues", middlewares.FCMInsertValues)
	a.DB.Prepare("FSCMSelectForumBySlug", middlewares.FSCMSelectForumBySlug)
	a.DB.Prepare("FSCMInsertValues", middlewares.FSCMInsertValues)
	a.DB.Prepare("TCMUpdateForumThreadsCount", middlewares.TCMUpdateForumThreadsCount)
	a.DB.Prepare("FSDGetValues", middlewares.FSDGetValues)
	a.DB.Prepare("TCMInsertToNewTable", middlewares.TCMInsertToNewTable)

}

//CloseConnection - closes database connection
func (a *Application) CloseConnection() {
	a.DB.Close()
}
