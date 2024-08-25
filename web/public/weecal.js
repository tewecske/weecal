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
