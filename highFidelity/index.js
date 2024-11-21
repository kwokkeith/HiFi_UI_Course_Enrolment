// Filter active state
const pillarFilter = document.getElementById("pillarFilter");
const pillarFilterInputList = pillarFilter.querySelectorAll("input");
pillarFilterInputList.forEach( input => {
    input.addEventListener("click", () => {
        if(!input.classList.contains("filterActiveState")){
            input.classList.add("filterActiveState");
        }
        else{
            input.classList.remove("filterActiveState");
        }
    })
});
const termFilter = document.getElementById("termFilter");
const termFilterInputList = termFilter.querySelectorAll("input");
termFilterInputList.forEach( input => {
    input.addEventListener("click", () => {
        if(!input.classList.contains("filterActiveState")){
            input.classList.add("filterActiveState");
        }
        else{
            input.classList.remove("filterActiveState");
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