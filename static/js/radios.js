$(document).ready(function () {
    $('input[type=radio]').change(function() {
        // When any radio button on the page is selected,
        // then deselect all other radio buttons.
        $('input[type=radio]:checked').not(this).prop('checked', false);
    });
});
