/// <reference path="./jquery-3.6.0.min.js"/>
function peRandID(length) {
    var result           = '';
    var characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
    var charactersLength = characters.length;
    for ( var i = 0; i < length; i++ ) {
        result = result + characters.charAt(Math.floor(Math.random() * charactersLength));
        // console.log(result);
    }
    return result
}
class myTextInput {
    constructor(params) {
        let thisObject = this
        params.domId = peRandID(15);
        thisObject.params = params;
        thisObject.jQUi = {};
        thisObject.jQContainer = jQuery(`#${params.containerId}`)
        thisObject.htmlTag = `
            <div id="${params.domId}" class="textInput">
                <div class="inputLabel">${params.label}: </div>
                <input id="${params.domId}_inb" type="text" class="inputBox"/>
                <div class="icon"></div>
                <div class="btn-container">
                    <div id="${params.domId}_btn" class="inputButton">Show Value</div>
                </div>
            </div>
        `;
        thisObject.attachUI().then(()=>{
            thisObject.jQUi.btn.click(()=>{
                thisObject.currentValue = thisObject.jQUi.inb.val();
                alert(thisObject.currentValue);
            })
            
            thisObject.jQUi.inb.keyup(()=>{
                console.clear();
                console.log(thisObject.jQUi.inb.val());
            });
        })
    }
    attachUI() {
        let thisObject = this;
        return new Promise((resolve, reject)=>{
            thisObject.jQContainer.append(thisObject.htmlTag);
            thisObject.jQUi.btn = jQuery(`#${thisObject.params.domId}_btn`);
            thisObject.jQUi.inb = jQuery(`#${thisObject.params.domId}_inb`);
            resolve();
        });        
    }
}