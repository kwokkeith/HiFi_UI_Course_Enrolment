try {
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

} catch (error) {
    console.error(error);
}


