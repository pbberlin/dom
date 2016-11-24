// ==UserScript==
// @name         citeproof-connector
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version.
// @description  http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.123
// @author       iche
// @downloadURL  https://raw.githubusercontent.com/pbberlin/dom/master/ui/tamper-monkey-citeproof-connector.js
// @updateURL    https://raw.githubusercontent.com/pbberlin/dom/master/ui/tamper-monkey-citeproof-connector.js
// @match        *://*/*
// @match        *://localhost:*/*
// @require      https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js
// @require      https://code.jquery.com/ui/1.12.1/jquery-ui.js
// @grant        none
// @noframes
// @run-at      document-end
// ==/UserScript==


// fallback http://encosia.com/3-reasons-why-you-should-let-google-host-jquery-for-you/
if (typeof jQuery === 'undefined') {
    console.log("CDN blocked by Iran or China?");
    document.write(unescape('%3Cscript%20src%3D%22/path/to/your/scripts/jquery-2.1.4.min.js%22%3E%3C/script%3E'));
}

(function ($, undefined) {
    $(function () {
        console.log("isolated jQuery start");


        function TriggerIt(){
            // This only works, if second argument matches opener.window.href
            // opener.postMessage("url-info", "http://localhost:8080");
            var data = {type: "url-info", info: { loc: window.location.href, remoteUrl: "aaa" }};


            // http://stackoverflow.com/questions/2120060
            var content = document.body.parentNode.innerHTML;
            var doctypeInformation = document.body.parentNode.previousSibling;
            data.info.html = content;
            opener.postMessage(data, "*");
            console.log("tamper monkey: window location is ",window.location.href);
        }



        function AddCssAndHtml(){
            if ($('#css-hover-popup').length <= 0) {
                var s =  '';
                s += '<style type="text/css"  id="css-hover-popup" >';
                s += '.ui-draggable-handle {';
                s += '    -ms-touch-action: none;';
                s += '    touch-action: none;';
                s += '}';
                s += '#id32168 { ';
                s += '    right: 10px; top: 10px; width: 150px; min-height: 40px;';
                s += '    position: absolute; z-index: 2100;';
                s += '    background-color: #aca;';
                s += '    border: 1px solid #464;';
                s += '    border-radius: 8px;';
                s += '    font-family: Monospace;';
                s += '    padding: 16px 14px;';
                s += '    font-size: 13px;';
                s += '    text-align: center;';
                s += '}';
                s += '</style>';
                $(s).appendTo('head');



                var cnt = "https://citeproof.appspot.com/upload-receiver";
                var popupUpScaffold = "<div id='bracket32168' style='position: relative;'><div id='id32168' >Click to save as Quotation. <span style='font-size:10px;'><br><br>or drag to move<span></div></div>";
                $('body').prepend(popupUpScaffold);  //next after <body;  most counter intuitive: stackoverflow.com/questions/5073016
            }
        }


        $( document ).ready(function() {
            console.log( "tamper monkey: document ready start" );


            console.log("tamper monkey: window name is",window.name);
            if (window.name==="browser_bridge_window") {
                console.log("active");
                AddCssAndHtml();
                $( function() {
                    $('#id32168').draggable({ cursor:"move" });
                    //$('#id32168').onclick = TriggerIt;
                    $('#id32168').click(TriggerIt);
                });

                $( function() {
                    //TriggerIt();
                });


            }




            console.log( "tamper monkey: document ready end" );
        });


        console.log("isolated jQuery end");
    });
})(window.jQuery.noConflict(true));