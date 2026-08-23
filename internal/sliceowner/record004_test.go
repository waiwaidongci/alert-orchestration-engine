package r004
import ("context"; "testing"; "time"; "github.com/example/alert-orchestration-engine/internal/domain/notification"; "github.com/example/alert-orchestration-engine/internal/domain/rule"; "github.com/example/alert-orchestration-engine/internal/infrastructure/queue")
func TestR004Batch(t *testing.T){x:=[]notification.Record{notification.New("n","a","r","webhook","x",time.Now())};b:=notification.NewBatch(x,time.Now());x[0].Target="changed";if b.Items[0].Target=="changed"{t.Fatal("batch alias")}}
func TestR004Sort(t *testing.T){x:=[]rule.Rule{{Name:"late",EscalationMinutes:2},{Name:"early",EscalationMinutes:1}};_=rule.Sort(x);if x[0].Name!="late"{t.Fatal("sort alias")}}
func TestR004Routes(t *testing.T){x:=[]notification.Route{{Channel:"webhook",Enabled:true,Target:"x"}};y:=notification.SelectRoutes(x,nil);x[0].Channel="changed";if y[0].Channel=="changed"{t.Fatal("route alias")}}
func TestR004Body(t *testing.T){q:=queue.New();b:=[]byte("ok");_=q.Publish(context.Background(),queue.Message{Topic:"t",Body:b});b[0]='X';if string(q.Messages[0].Body)!="ok"{t.Fatal("body alias")}}
