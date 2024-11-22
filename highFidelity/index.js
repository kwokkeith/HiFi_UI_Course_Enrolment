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
const selectAllCheckBox = document.getElementById("shoppingCartSelectAll");
selectAllCheckBox.addEventListener("change", () => {
    const checkboxes = document.querySelectorAll(".shoppingCartCourseBlockCheckbox");
    checkboxes.forEach( checkbox => {
        checkbox.checked = selectAllCheckBox.checked;
    });
});

// Course block active state (courseBlock button -> courseDetail)
// All courseDetail should contains ".hiddenState" by default
const courseFullContainerList = document.querySelectorAll(".courseFullContainer");
const courseBlockList = document.querySelectorAll(".courseBlock");
courseBlockList.forEach( block => {
    block.addEventListener("click", () => {

        // Hide all courseFullContainer except the one users clicked
        for(let i = 0; i < courseFullContainerList.length; i++){
            if(courseFullContainerList[i].children[0] !== block){
                courseFullContainerList[i].classList.add("hiddenState");
            }
            else{
                // Hide this courseBlcok 
                block.classList.add("hiddenState");
                // Show the courseDetail
                courseFullContainerList[i].children[1].classList.remove("hiddenState");
            }
        }
    })
});
// Return to course block page
const courseDetailReturnButtonList = document.querySelectorAll(".courseDetailReturnButton");
courseDetailReturnButtonList.forEach( button => {
    button.addEventListener("click", () => {
        // Hide courseDetail and show courseBlock for this course
        // button -> courseDetailHeader -> courseDetailSubContainer -> courseDetail
        const thisCourseDetail = button.parentElement.parentElement.parentElement;
        thisCourseDetail.classList.add("hiddenState");
        const thisCourseBlock = thisCourseDetail.previousElementSibling;
        thisCourseBlock.classList.remove("hiddenState");
        // Show all other courseBlock
        for(let i = 0; i < courseFullContainerList.length; i++){
            if(courseFullContainerList[i].children[0] !== thisCourseBlock){
                courseFullContainerList[i].classList.remove("hiddenState");
            }
        }
        resetCollapsibleButton();
    })
});

// collapsibleButton (Student Review Button)
const collapsibleButton = document.getElementsByClassName("collapsibleButton");
for(let button of collapsibleButton){
    button.addEventListener("click", () => {

        const collapsibleContent = button.nextElementSibling;

        // Ensure that it is a collapsibleContent
        if(collapsibleContent.classList.contains("collapsibleContent")){
            // State checking
            if(collapsibleContent.classList.contains("hiddenState")){
                collapsibleContent.classList.remove("hiddenState");
            }
            else{  
                collapsibleContent.classList.add("hiddenState");
            }
        }
    });
}
// Reset all collapsibleButton
function resetCollapsibleButton(){
    const collapsibleContentList = document.getElementsByClassName("collapsibleContent");
    for(let element of collapsibleContentList){
        // if the collapsible content is not in hidden state
        if(!element.classList.contains("hiddenState")){
            element.classList.add("hiddenState");
        }
    }
}

// Pop up window
const popuplist = document.getElementsByClassName("popup");
for(let popup of popuplist){
    popup.addEventListener("mouseover", () => {
        const popuptext = popup.children[0];
        popuptext.classList.toggle("showPopup");
    });
    popup.addEventListener("mouseout", () => {
        const popuptext = popup.children[0];
        popuptext.classList.remove("showPopup"); // Hide the popup
    });
}

