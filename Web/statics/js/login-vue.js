var app = new Vue({
    el : "#container",
    data : {
        userID : null,
        password : '',

        isRegister:false,
        // 图片路径请自行更改
        imgSrc : ["/img/1.jpg","/img/2.jpg","/img/3.jpg","/img/4.jpg",
    ],
        imgIndex : 0,
    },
    created: function () {
        setInterval(this.lantenSlide, 2000);
    },

    methods : {
        selectLogin:function(){
            // alert("login")
            this.isRegister = false;
        },
        selectRegister:function(){
            // alert("register")
            this.isRegister = true;
        }, 
        loginDo: function() {
            var that = this
            axios.post('login', {
                userID: this.userID,
                userPwd: this.password,
              })
              .then(function (response) {
                  console.log(response.data);
               if (response.data.res == "ok") {
                alert("登录成功");
                window.location.href = "./hall"
               } else {
                   alert("登陆失败");
               }
               
              })
              .catch(function (error) {
                console.log(error);
              });
        },

        lantenSlide : function(){
            this.imgIndex = ((this.imgIndex+1) % this.imgSrc.length)
            // console.log(this.imgIndex);
        },

    }
});   