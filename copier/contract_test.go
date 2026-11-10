package copier
import (
 "context"
 "testing"
 "google.golang.org/grpc"
 "google.golang.org/grpc/metadata"
 pb "github.com/MetaRPC/GoCopier/copier/copierpb"
)
type spy struct {t *testing.T; count int; start *pb.StartRequest}
func(s *spy) check(ctx context.Context){s.count++;md,_:=metadata.FromOutgoingContext(ctx);if md.Get("apikey")[0]!="owner-test"{s.t.Fatal("missing API key")};if _,ok:=ctx.Deadline();!ok{s.t.Fatal("no deadline")}}
func(s *spy) Start(c context.Context,r *pb.StartRequest,o ...grpc.CallOption)(*pb.StartReply,error){s.check(c);s.start=r;return &pb.StartReply{Ok:true},nil}
func(s *spy) List(c context.Context,r *pb.ListRequest,o ...grpc.CallOption)(*pb.ListReply,error){s.check(c);return &pb.ListReply{Ok:false,Error:"refused"},nil}
func(s *spy) Pause(c context.Context,r *pb.PauseRequest,o ...grpc.CallOption)(*pb.SimpleReply,error){s.check(c);if r.Paused{s.t.Fatal("resume changed")};return &pb.SimpleReply{Ok:true},nil}
func(s *spy) Remove(c context.Context,r *pb.RemoveRequest,o ...grpc.CallOption)(*pb.SimpleReply,error){s.check(c);return &pb.SimpleReply{Ok:true},nil}
func TestAllMethods(t *testing.T){
 mock:=&spy{t:t};a:=&CopierAccount{UserKey:"owner-test",ManagerKey:"manager-test"};s:=&CopierService{Account:a,grpcClient:mock}
 c:=metadata.NewOutgoingContext(context.Background(),metadata.Pairs("trace-id","preserved"))
 if md,_:=metadata.FromOutgoingContext(a.AuthContext(c));md.Get("trace-id")[0]!="preserved"{t.Fatal("discarded caller metadata")}
 s.Start(c,&StartRequest{Master:&Account{},Slave:&Account{},CopySl:false,CopyTp:true,RiskType:"FixedLot",RiskValue:"0.01"})
 result,_:=s.List(c);if result.Ok||result.Error!="refused"{t.Fatal("lost failure")}
 s.Pause(c,"test",false);s.Remove(c,"test")
 if mock.count!=4||mock.start.UserKey!="owner-test"||mock.start.ManagerKey!="manager-test"||mock.start.CopySl||!mock.start.CopyTp{t.Fatal("request changed")}
 if _,err:=s.Start(c,nil);err==nil{t.Fatal("nil request accepted")}
}
func TestInvalidKeys(t *testing.T){for _,k:=range []string{"","TRIAL","trial","a b","x\r\nInjected: yes"}{if _,err:=NewCopierAccount("copy.mrpc.pro",k);err==nil{t.Fatalf("accepted %q",k)}}}
