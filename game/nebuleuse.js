 APIURL = "http://orypie:8080"
 Vue.http.interceptors.push(function(request, next) {
          // modify headers
      request.headers.set('accept', 'application/json');
    
      // continue to next interceptor
      next();
    });

 var app = new Vue({
  el: '#app',
  data: {
    message: 'Hello Vue!',
    sessionid: ""
  },
  created: function () {
    // `this` points to the vm instance
    
    this.$http.get(APIURL + '/status').then(response => {
        
            // get body data
            //this.message = response.body;

            // Login
            this.Login();
        
          }, response => {
            // error callback
          });
  },
  methods:{
      Login: function(){
        console.log('hello');
        
        var formData = new FormData();
        formData.append('username', 'test2');
        formData.append('password', 'test');

        this.$http.post(APIURL + '/connect', formData).then(response => {
                this.message = response.body;
                this.sessionid = response.body.SessionId;
                //this.SubscribeToGame(response.body.SessionId);
                
                this.LongPoll(response.body.SessionId);
                
              }, response => {
                  this.message = response;
              });
      },
      LongPoll: function(SessionId){

        var formData = new FormData();
        formData.append('sessionid', SessionId);
        this.$http.post(APIURL + '/getMessages', formData).then(response => {
                this.message = response.body;
                this.LongPoll(SessionId);
              }, response => {
                  this.message = response;
                  this.LongPoll(SessionId);
              });
      },
      SendMessage: function(msg){
        console.log(this.sessionid);
        var formData = new FormData();
        formData.append('sessionid', this.sessionid);
        formData.append('pipe', "game");
        formData.append('channel', "players");
        formData.append('message', msg);
        this.$http.post(APIURL + '/sendMessage', formData).then(response => {
                this.message = response.body;
              }, response => {
                  this.message = response;
              });

      },

      SubscribeToGame: function(sessionid){
        var formData = new FormData();
        console.log(sessionid);
        formData.append('sessionid', sessionid);
        formData.append('pipe', "game");
        formData.append('channel', "players");
        this.$http.post(APIURL + '/subscribeTo', formData).then(response => {
                console.log(response);
              }, response => {
                console.log(response);
              });

      }
  }
  
})