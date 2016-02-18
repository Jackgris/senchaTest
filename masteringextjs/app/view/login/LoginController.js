Ext.define('Packt.view.login.LoginController', {
    extend: 'Ext.app.ViewController',
    alias: 'controller.login',

    requires:[
        'Packt.view.login.CapsLockTooltip',
        'Packt.util.Util'
    ],
    
    onTextFieldSpecialKey: function(field, e, options){
        if (e.getKey() === e.ENTER) {
            this.doLogin();
        }
    },

    onTextFieldKeyPress: function(field, e, options){
        var charCode = e.getCharCode(),
            me = this;

        if ((e.shiftKey && charCode >= 97 && charCode <= 122) ||
            (!e.shiftKey && charCode >= 65 && charCode <= 90)){

            if(me.capslockTooltip === undefined){
                me.capslockTooltip = Ext.widget('capslocktooltip');
            }

            me.capslockTooltip.show();
        } else {
            if(me.capslockTooltip !== undefined){
                me.capslockTooltip.hide();
            }
        }
    },

    onButtonClickCancel: function(button, e, options){
        this.lookupReference('form').reset();
    },

    onButtonClickSubmit: function(button, e, options){
        var me = this;

        if(me.lookupReference('form').isValid()){
            me.doLogin();
        }
    },

    doLogin: function(){
        var me = this,
            form = me.lookupReference('form');
        this.getView().mask('Authenticating...Please wait...');
        form.submit({
            clientValidation: true,
            url: 'security/signup',
            scope: me,
            success: 'onLoginSuccess',
            failure: 'onLoginFailure'
        });
    },

    onLoginFailure: function(form, action){
        this.getView().unmask();
        var result = Packt.util.Util.decodeJSON(action.response.responseText);

        console.log("Llego")
        console.log(action.response.responseText)
        console.log(result)
        console.log(action.response.status)
        console.log(action.response.statusText)
        
        switch(action.failureType){
        case Ext.form.action.Action.CLIENT_INVALID:
            Packt.util.Util.showErrorMsg(
                'Form fields may not be submitted with invalid values, client invalid');
            break;
        case Ext.form.action.Action.CONNECT_FAILURE:
            Packt.util.Util.showErrorMsg(
                'Form fields may not be submitted with invalid values, connect failure');
            break;
        case Ext.form.action.Action.SERVER_INVALID:
            Packt.util.Util.showErrorMsg(result.msg);
        }
    },

    onLoginSuccess: function(form, action){
        this.getView().unmask();
        this.getView().close();
        Ext.create('Packt.view.main.Main');
    }
});
