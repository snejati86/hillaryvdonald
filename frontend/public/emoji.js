const MAX_LENGTH = 38;
const EMOJI_LENGTH = 2;
const MARKER = '.';

function appendEmoji(emoji,field){
    //console.log(field.text().length)
    if ( field.text() === MARKER){
        //FIRST ONE.
        field.text((field.text().substr(MARKER.length)));
    }
    else{
        if ( field.text().length > MAX_LENGTH  ){
            field.text(field.text().substr(EMOJI_LENGTH));

        }
    }
    field.append(emoji);
}

function displayTweet ( tweet , el  ){
    if( !el.is(':animated') ){
        el[0].innerHTML= tweet;


        el.animate({opacity: 1}, 1000, 'linear', function() {
            el.animate({opacity: 0}, 2500, 'linear');
        })
    }
}
function flash(logo){
    if( !logo.is(':animated') ) {
        logo.animate({opacity: 0.3}, 300, 'swing', function () {
            logo.animate({opacity: 0.8}, 300, 'swing');
        });
    }
}
