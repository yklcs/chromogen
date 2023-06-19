const $ = document.querySelector.bind(document)
const $$ = document.querySelectorAll.bind(document)

$("button#viewmode-toggler").addEventListener("click", (ev) => {
  const button = ev.target
  const thumbs = $("div#thumbs")

  const gridText = "Grid"
  const galleryText = "Gallery"

  thumbs.classList.toggle("thumbs-gallery")
  thumbs.classList.toggle("thumbs-grid")

  if (button.innerHTML === galleryText) {
    button.innerHTML = gridText
  } else if (button.innerHTML === gridText) {
    button.innerHTML = galleryText
  }
})

$$("figure").forEach((fig) => {
  const placeholder = fig.querySelector("img.placeholder")
  const thumb = fig.querySelector("img.thumb")
  if (!thumb.complete) {
    thumb.addEventListener("load", (ev) => {
      thumb.classList.remove("thumb-unloaded")
      placeholder.classList.remove("placeholder-unloaded")
    })
  } else {
    thumb.classList.remove("thumb-unloaded")
    placeholder.classList.remove("placeholder-unloaded")
  }
})
