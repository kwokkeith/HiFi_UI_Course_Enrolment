var pillars = {"epd": 0,
               "esd": 0,
               "istd": 0,
               "asd": 0,
               "dai": 0,
               "hass": 0};

var terms = {"1": 0,
             "2": 0,
             "3": 0,
             "4": 0,
             "5": 0,
             "6": 0,
             "7": 0,
             "8": 0};

var nopillars = {"epd": 0,
               "esd": 0,
               "istd": 0,
               "asd": 0,
               "dai": 0,
               "hass": 0};

var noterms = {"1": 0,
             "2": 0,
             "3": 0,
             "4": 0,
             "5": 0,
             "6": 0,
             "7": 0,
             "8": 0};

function resetPillars() {
    const pillarFilter = document.getElementById("pillarFilter");
    const pillarFilterInputList = pillarFilter.querySelectorAll("input");
    pillarFilterInputList.forEach( input => {
            if(input.classList.contains("filterActiveState")){
                input.classList.remove("filterActiveState");
                pillars[input.value.toLowerCase()] = 0;
            }
    });
}

function resetTerms() {
    const termFilter = document.getElementById("termFilter");
    const termFilterInputList = termFilter.querySelectorAll("input");
    termFilterInputList.forEach( input => {
            if(input.classList.contains("filterActiveState")){
                input.classList.remove("filterActiveState");
                terms[input.value] = 0;
            }
    });
}

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
        showCourse();
    });
});

function toggleCheckbox(n) {
    const checkboxes = document.querySelectorAll(".shoppingCartCourseBlockCheckbox");
    checkboxes[n].checked ^= 1;
    showCourse();
}

// Shopping Cart/Enroll stuff
var cartCount = 0;

const scheduleContainerCourses = ['<div class="course1fade">20.420<br>CI01</div>',
    '<div class="course2fade">20.420<br>CI01</div>',
    '<div class="course3fade">40.008<br>CI02</div>',
    '<div class="course4fade">40.008<br>CI02</div>'];

const cartCourses = [' <div onClick="toggleCheckbox(0)" class="shoppingCartCourseBlock"> <input onClick="showCourse(); toggleCheckbox(0);" id="shoppingCartCourseBlockCheckbox1" type="checkbox" class="shoppingCartCourseBlockCheckbox"> <div class="shoppingCartCourseBlockCourse"> <span class="shoppingCartCourseBlockCourseCode">20.420</span> <span class="shoppingCartCourseBlockCourseTitle">Green Architecture and Urban Sustainability</span> </div> <button type="button" class="shoppingCartCourseBlockRemoveButton"> <img onClick="removeCourse(0)" src="./images/shoppingCart/cancelIcon.svg"> </button> </div> ',' <div onClick="toggleCheckbox(1)" class="shoppingCartCourseBlock"> <input onClick="showCourse(); toggleCheckbox(1);" id="shoppingCartCourseBlockCheckbox2" type="checkbox" class="shoppingCartCourseBlockCheckbox"> <div class="shoppingCartCourseBlockCourse"> <span class="shoppingCartCourseBlockCourseCode">40.008</span> <span class="shoppingCartCourseBlockCourseTitle">Systems Thinking for Operational Excellence</span> </div> <button type="button" class="shoppingCartCourseBlockRemoveButton"> <img onClick="removeCourse(1)" src="./images/shoppingCart/cancelIcon.svg"> </button> </div> '];

function showCourse() {
    document.getElementById("scheduleContainer").innerHTML = '<img class="scheduleImage" src="./images/schedule/emptySchedule.png">';
    if (document.getElementById("shoppingCartCourseBlockCheckbox1").checked && !document.getElementById("shoppingCartCourseBlockCheckbox2").checked) {
        document.getElementById("scheduleContainer").innerHTML += scheduleContainerCourses[0] + scheduleContainerCourses[1];
    }
    if (!document.getElementById("shoppingCartCourseBlockCheckbox1").checked && document.getElementById("shoppingCartCourseBlockCheckbox2").checked) {
        document.getElementById("scheduleContainer").innerHTML += scheduleContainerCourses[2] + scheduleContainerCourses[3];
    }
    if (document.getElementById("shoppingCartCourseBlockCheckbox1").checked && document.getElementById("shoppingCartCourseBlockCheckbox2").checked) {
        document.getElementById("scheduleContainer").innerHTML += scheduleContainerCourses[0] + scheduleContainerCourses[1];
        document.getElementById("scheduleContainer").innerHTML += scheduleContainerCourses[2] + scheduleContainerCourses[3];
    }
}

function removeCourse(n) {
    if (cartCount == 2) { 
        document.getElementById("shoppingCartWindow").innerHTML = cartCourses[(n+1)%2];
        cartCount -= 1;
    }
    else if (cartCount == 1) { 
        document.getElementById("shoppingCartWindow").innerHTML = "";
        cartCount -= 1;
    }
}

function addCourse() {
    if (cartCount < 2) { 
        document.getElementById("shoppingCartWindow").innerHTML += cartCourses[cartCount];
        cartCount += 1;
    }
}

function enrollCourse(code, name) {
    document.getElementById("scheduleContainer").innerHTML = '<img class="scheduleImage" src="./images/schedule/emptySchedule.png"> <div class="course1">'+code+'<br>CI01</div> <div class="course2">'+code+'<br>CI01</div>';
    alert("Successfully Enrolled into " + code + ": " + name + "!");
}

function enrollAll() {
    document.getElementById("scheduleContainer").innerHTML = '<img class="scheduleImage" src="./images/schedule/emptySchedule.png">';
    if (document.getElementById("shoppingCartCourseBlockCheckbox1").checked && !document.getElementById("shoppingCartCourseBlockCheckbox2").checked) { 
        document.getElementById("scheduleContainer").innerHTML += '<div class="course1">20.420<br>CI01</div> <div class="course2">20.420<br>CI01</div>';
        document.getElementById("shoppingCartWindow").innerHTML = cartCourses[1];
        alert("Successfully Enrolled!");
    }
    if (!document.getElementById("shoppingCartCourseBlockCheckbox1").checked && document.getElementById("shoppingCartCourseBlockCheckbox2").checked) { 
        document.getElementById("scheduleContainer").innerHTML += '<div class="course3">40.008<br>CI02</div> <div class="course4">40.008<br>CI02</div>';
        document.getElementById("shoppingCartWindow").innerHTML = cartCourses[0];
        alert("Successfully Enrolled!");
    }
    if (document.getElementById("shoppingCartCourseBlockCheckbox1").checked && document.getElementById("shoppingCartCourseBlockCheckbox2").checked) { 
        document.getElementById("scheduleContainer").innerHTML += '<div class="course1">20.420<br>CI01</div> <div class="course2">20.420<br>CI01</div>';
        document.getElementById("scheduleContainer").innerHTML += '<div class="course3">40.008<br>CI02</div> <div class="course4">40.008<br>CI02</div>';
        document.getElementById("shoppingCartWindow").innerHTML = '';
        alert("Successfully Enrolled!");
    }
}

function clearAll() {
    cartCount = 0;
    document.getElementById("scheduleContainer").innerHTML = '<img class="scheduleImage" src="./images/schedule/emptySchedule.png">';
    document.getElementById("shoppingCartWindow").innerHTML = '';
}


