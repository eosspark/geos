function createaccount(input) {
    var creator=input['creator']
    var  name =input['name']
    var owner =input['owner']
    var active = input['active']
    console.log(creator,name,owner,active)
    eos.CreateAccount(creator,name,owner,active)

}

function pushaction(input){
    var account = input['account']
    var action = input['action']
    var data = input['data']
    var permisson = input['permission']
    console.log(account,action,data,permisson)
    eos.PushAction(account,action,data,permisson)
}


// function getTable(input){
//
// }




