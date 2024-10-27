function dropDragged(event) {
  const target = event.target.closest('td');
  const element = document.getElementById(event.dataTransfer.getData('text/plain'));
  target.appendChild(element);
}

document.addEventListener('alpine:init', () => {
  Alpine.data('draggable', () => ({
    dragging: false,
    // init() {
    //   console.log("alpineDraggable initialized");
    // },
    dragStart(event) {
      this.dragging = true;
      event.dataTransfer.effectAllowed='move';
      event.dataTransfer.setData('text/plain', event.target.id);
    },
    dragEnd(event) {
      this.dragging = false;
    }
  }))
})

htmx.on("htmx:load", function(evt) {
  if (htmx.find("#teams-container") !== null) {
    htmx.on('#teams-container', 'htmx:afterSwap', function(event) {
        if (event.detail.target.id === 'viewTeamModal') {
            viewTeamModal.showModal();
        }
    });
  }
});

async function confirmation(e) {
  const messageElement = document.querySelector("#deleteModal #confirmDeleteMessage");
  messageElement.innerHTML = e.detail.question;
  const confirmDeleteButton = document.querySelector("#deleteModal #confirmDeleteButton");
  const cancelDeleteButton = document.querySelector("#deleteModal #cancelDeleteButton");

  const promise = await new Promise((resolve, reject) => {
    deleteModal.showModal();

    confirmDeleteButton.addEventListener("click", () => {
      resolve(true);
    });

    cancelDeleteButton.addEventListener("click", () => {
      resolve(false);
    });
  });

  if (promise) {
    e.detail.issueRequest(true);
  }
}

document.addEventListener("htmx:confirm", function(e) {
  if (!e.detail.target.hasAttribute('hx-confirm')) return

  e.preventDefault();

  confirmation(e);

});

