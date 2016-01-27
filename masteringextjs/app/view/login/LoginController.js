Ext.define('Packt.view.login.LoginController', {
    extend: 'Ext.app.ViewController',
    alias: 'controller.login',

    onTextFieldSpecialKey: function(field, e, options){},

    onTextFieldKeyPress: function(field, e, options){},

    onButtonClickCancel: function(button, e, options){
        console.log('login cancel');
        this.lookupReference('form').reset();
    },

    onButtonClickSubmit: function(button, e, options){
        console.log('login submit');
        var me = this;
        if(me.lookupReference('form').isValid()){
            me.doLogin();
        }
    },

    doLogin: function(){
        var me = this,
            form = me.lookupReference('form');
        form.submit({
            clientValidation: true,
            url: 'security/signup',
            scope: me,
            success: 'onLoginSuccess',
            failure: 'onLoginFailure'
        });
    },

    onLoginFailure: function(form, action){},

    onLoginSuccess: function(form, action){}
});
