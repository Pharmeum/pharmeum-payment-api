package db

import (
	"context"

	"github.com/go-kivik/kivik"
	"github.com/go-kivik/kiviktest/kt"
)

func init() {
	kt.Register("DeleteAttachment", delAttachment)
}

func delAttachment(ctx *kt.Context) {
	ctx.RunRW(func(ctx *kt.Context) {
		dbname := ctx.TestDB()
		defer ctx.DestroyDB(dbname)
		ctx.Run("group", func(ctx *kt.Context) {
			ctx.RunAdmin(func(ctx *kt.Context) {
				ctx.Parallel()
				testDeleteAttachments(ctx, ctx.Admin, dbname, "foo.txt")
				testDeleteAttachments(ctx, ctx.Admin, dbname, "NotFound")
				testDeleteAttachmentsDDoc(ctx, ctx.Admin, dbname, "foo.txt")
				testDeleteAttachmentNoDoc(ctx, ctx.Admin, dbname)
			})
			ctx.RunNoAuth(func(ctx *kt.Context) {
				ctx.Parallel()
				testDeleteAttachments(ctx, ctx.NoAuth, dbname, "foo.txt")
				testDeleteAttachments(ctx, ctx.NoAuth, dbname, "NotFound")
				testDeleteAttachmentsDDoc(ctx, ctx.NoAuth, dbname, "foo.txt")
				testDeleteAttachmentNoDoc(ctx, ctx.NoAuth, dbname)
			})
		})
	})
}

func testDeleteAttachmentNoDoc(ctx *kt.Context, client *kivik.Client, dbname string) {
	db := client.DB(context.Background(), dbname, ctx.Options("db"))
	if err := db.Err(); err != nil {
		ctx.Fatalf("Failed to connect to db")
	}
	ctx.Run("NoDoc", func(ctx *kt.Context) {
		ctx.Parallel()
		_, err := db.DeleteAttachment(context.Background(), "nonexistantdoc", "2-4259cd84694a6345d6c534ed65f1b30b", "foo.txt")
		ctx.CheckError(err)
	})
}

func testDeleteAttachments(ctx *kt.Context, client *kivik.Client, dbname, filename string) {
	ctx.Run(filename, func(ctx *kt.Context) {
		doDeleteAttachmentTest(ctx, client, dbname, ctx.TestDBName(), filename)
	})
}

func testDeleteAttachmentsDDoc(ctx *kt.Context, client *kivik.Client, dbname, filename string) {
	ctx.Run("DesignDoc/"+filename, func(ctx *kt.Context) {
		doDeleteAttachmentTest(ctx, client, dbname, "_design/"+ctx.TestDBName(), filename)
	})
}

func doDeleteAttachmentTest(ctx *kt.Context, client *kivik.Client, dbname, docID, filename string) {
	db := client.DB(context.Background(), dbname, ctx.Options("db"))
	if err := db.Err(); err != nil {
		ctx.Fatalf("Failed to connect to db")
	}
	ctx.Parallel()
	adb := ctx.Admin.DB(context.Background(), dbname, ctx.Options("db"))
	if err := adb.Err(); err != nil {
		ctx.Fatalf("Failed to open db: %s", err)
	}
	doc := map[string]interface{}{
		"_id": docID,
		"_attachments": map[string]interface{}{
			"foo.txt": map[string]interface{}{
				"content_type": "text/plain",
				"data":         "VGhpcyBpcyBhIGJhc2U2NCBlbmNvZGVkIHRleHQ=",
			},
		},
	}
	rev, err := adb.Put(context.Background(), docID, doc)
	if err != nil {
		ctx.Fatalf("Failed to create doc: %s", err)
	}
	rev, err = db.DeleteAttachment(context.Background(), docID, rev, filename)
	if !ctx.IsExpectedSuccess(err) {
		return
	}
	var i interface{}
	if err = db.Get(context.Background(), docID, map[string]interface{}{"rev": rev}).ScanDoc(&i); err != nil {
		ctx.Fatalf("Failed to get deleted doc: %s", err)
	}
}
