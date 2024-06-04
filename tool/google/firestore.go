package google

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
)

type Firestore struct {
	*firestore.Client
}

func NewFirestore(ctx context.Context, app *firebase.App) (cli *Firestore, err error) {
	cli = &Firestore{}
	cli.Client, err = app.Firestore(ctx)

	return
}
