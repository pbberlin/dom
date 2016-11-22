// ==UserScript==
// @name         citeproof-connector
// @description  include jQuery and make sure window.$ is the content page's jQuery version, and this.$ is our jQuery version.
// @description  http://stackoverflow.com/questions/28264871/require-jquery-to-a-safe-variable-in-tampermonkey-script-and-console
// @namespace    http://your.homepage/
// @version      0.123
// @author       iche
// @downloadURL
// @updateURL
// @match        *://*.welt.de/*
// @match        *://citeproof.appspot.com/*
// @match        *://*.economist.com/*
// @require      https://ajax.googleapis.com/ajax/libs/jquery/2.1.4/jquery.min.js
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

$( document ).ready(function() {

    console.log( "tamper monkey document ready start" );

    var hidingCountdownID = 0;

    console.log("window name is",window.name);
    if (window.name==="research") {
        console.log("active");
    }


    $( 'a' ).on( "mouseenter", function(evt) {
    });

    $( 'a' ).on( "mouseleave", function(evt) {
    });


    $( '#popup1' ).on( "mouseenter", function(evt) {
    });

    $( '#popup1' ).on( "mouseleave", function(evt) {
    });



    $( 'a' ).on( "focusin", function(evt) {
    });

    $( '#popup1' ).on( "focusout", function(evt) {
    });

    console.log( "tamper monkey document ready end" );
});

        console.log("isolated jQuery end");
    });
})(window.jQuery.noConflict(true));