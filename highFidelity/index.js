var pillars = {"epd": 0,
               "esd": 0,
               "istd": 0,
               "asd": 0,
               "dai": 0,
               "freshmore": 0};

var terms = {"1": 0,
             "2": 0,
             "3": 0,
             "4": 0,
             "5": 0,
             "6": 0,
             "7": 0,
             "8": 0};

// Filter active state
const pillarFilter = document.getElementById("pillarFilter");
const pillarFilterInputList = pillarFilter.querySelectorAll("input");
pillarFilterInputList.forEach( input => {
    input.addEventListener("click", () => {
        if(!input.classList.contains("filterActiveState")){
            input.classList.add("filterActiveState");
            pillars[input.value.toLowerCase()] ^= 1;
        }
        else {
            input.classList.remove("filterActiveState");
            pillars[input.value.toLowerCase()] ^= 1;
        }
    })
});
const termFilter = document.getElementById("termFilter");
const termFilterInputList = termFilter.querySelectorAll("input");
termFilterInputList.forEach( input => {
    input.addEventListener("click", () => {
        if(!input.classList.contains("filterActiveState")){
            input.classList.add("filterActiveState");
            terms[input.value] ^= 1;
        }
        else {
            input.classList.remove("filterActiveState");
            terms[input.value] ^= 1;
        }
    })
});

// Select all check box in shopping cart
selectAllCheckBox = document.getElementById("shoppingCartSelectAll")
selectAllCheckBox.addEventListener("change", () => {
    const checkboxes = document.querySelectorAll(".shoppingCartCourseBlockCheckbox");
    checkboxes.forEach( checkbox => {
        checkbox.checked = selectAllCheckBox.checked;
    });
});
