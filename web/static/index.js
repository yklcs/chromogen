/// <reference lib="dom" />

/**  @type {typeof document.querySelector} */
const $ = document.querySelector.bind(document)
/**  @type {typeof document.querySelectorAll} */
const $$ = document.querySelectorAll.bind(document)

const viewmodes = ["gallery", "grid"]

const initialViewmode = localStorage.getItem("viewmode")
if (initialViewmode !== null) {
  $("button#viewmode-toggler").innerHTML = initialViewmode
  $("div#thumbs").className = ""
  $("div#thumbs").classList.add(`thumbs-${initialViewmode}`)
}

$("button#viewmode-toggler").classList.toggle("noscript")
$("button#viewmode-toggler").addEventListener("click", (ev) => {
  const button = ev.target

  const idx = viewmodes.findIndex((mode) => mode === button.innerHTML)
  const newmode = viewmodes[(idx + 1) % viewmodes.length]
  localStorage.setItem("viewmode", newmode)
  button.innerHTML = newmode
  $("div#thumbs").className = ""
  $("div#thumbs").classList.add(`thumbs-${newmode}`)
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
