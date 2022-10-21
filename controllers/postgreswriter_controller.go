/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	"github.com/WANNA959/postgres-writer-operator/pkg/postgre"
	"github.com/go-logr/logr"

	dbv1 "github.com/WANNA959/postgres-writer-operator/api/v1"
	"github.com/WANNA959/postgres-writer-operator/pkg/model"
)

// PostgresWriterReconciler reconciles a PostgresWriter object
type PostgresWriterReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	PostgreClient *postgre.PostgresDBClient
}

var (
	logger     logr.Logger
	finalizers []string = []string{"finalizers.postgreswriters.demo.yash.com/cleanup-row"}
)

//+kubebuilder:rbac:groups=db.godx.com,resources=postgreswriters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=db.godx.com,resources=postgreswriters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=db.godx.com,resources=postgreswriters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PostgresWriter object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.12.2/pkg/reconcile
func (r *PostgresWriterReconciler) Reconcile(ctx context.Context, req ctrl.Request) (reconcile.Result, error) {
	logger = log.FromContext(ctx)
	// TODO(user): your logic here

	// parsing the incoming postgres-writer resource
	// 注意：Controller 逻辑是会在 CRD 状态变化的时候回调这个方法,此处需要自己判断这个事件是create/update/delete
	// 需要注意的是，如果没有其他辅助方式，你很难知道 CRD 是新建还是更新了
	// 但是，CRD 都对应着一个资源的状态，你只需要判断 CRD 的对象和资源的状态是否一致即可，至于是新建和更新，大多数情况都不是那么重要。
	postgreWriterObj := &dbv1.PostgresWriter{}
	err := r.Get(ctx, types.NamespacedName{Name: req.Name, Namespace: req.Namespace}, postgreWriterObj)
	logger.Info(fmt.Sprintf("get crd object:%+v\n", postgreWriterObj))

	if err != nil {
		if errors.IsNotFound(err) {
			// delete事件
			logger.Info(fmt.Sprintf("%s/%s has been delete", req.Namespace, req.Name))
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Error occurred while fetching the PostgresWriter resource")
		return ctrl.Result{}, err
	}

	// no delete, add finalizers
	if postgreWriterObj.GetDeletionTimestamp().IsZero() {
		if !reflect.DeepEqual(postgreWriterObj.GetFinalizers(), finalizers) {
			postgreWriterObj.SetFinalizers(finalizers)
			if err := r.Update(ctx, postgreWriterObj); err != nil {
				logger.Error(err, "error occurred while setting the finalizers of the PostgresWriter resource")
				return ctrl.Result{}, err
			}
		}
		if err := r.InsertOrUpdatePostgres(postgreWriterObj); err != nil {
			return ctrl.Result{}, err
		}
	}

	// delete
	if !postgreWriterObj.GetDeletionTimestamp().IsZero() {
		r.DeletePostgres(ctx, postgreWriterObj)
	}

	return ctrl.Result{}, nil
}

func (r *PostgresWriterReconciler) InsertOrUpdatePostgres(postgreWriterObj *dbv1.PostgresWriter) error {
	// parsing the table, name, age, country fields from the spec of the incoming postgres-writer resource
	id, name, age, department := postgreWriterObj.Spec.Id, postgreWriterObj.Spec.Name, postgreWriterObj.Spec.Age, postgreWriterObj.Spec.Department

	// forming a unique id corresponding to the incoming resource
	crdid := postgreWriterObj.Namespace + "/" + postgreWriterObj.Name
	logger.Info(fmt.Sprintf("InsertOrUpdatePostgres id: %+v, name: %+v, age: %+v, department: %+v, into table Student", id, name, age, department))
	// performing the `INSERT` to the DB with the provided name, age, country on the provided table
	student := &model.Student{
		Id:         id,
		Name:       name,
		Age:        age,
		Department: crdid,
	}
	std := &model.Student{}
	if err := std.Insert(student); err != nil {
		logger.Error(err, "error occurred while inserting the row in the Postgres DB")
		return err
	}
	return nil
}

func (r *PostgresWriterReconciler) DeletePostgres(ctx context.Context, postgreWriterObj *dbv1.PostgresWriter) error {
	// parse the id of the row to delete
	id := postgreWriterObj.Spec.Id
	std := &model.Student{}
	if err := std.Delete(id); err != nil {
		return fmt.Errorf("error occurred while running the DELETE query on Postgres: %w", err)
	}

	// remove the cleanup-row finalizer from the postgresWriterObject
	controllerutil.RemoveFinalizer(postgreWriterObj, finalizers[0])
	if err := r.Update(ctx, postgreWriterObj); err != nil {
		return fmt.Errorf("error occurred while running the DELETE query on Postgres: %w", err)
	}
	logger.Info("cleaned up the 'finalizers.postgreswriters.demo.yash.com/cleanup-row' finalizer successfully")
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostgresWriterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).For(&dbv1.PostgresWriter{}).Complete(r)
}
