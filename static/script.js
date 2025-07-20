const storageKey = "theme-preference";

const getColorPref = () => {
  if (localStorage.getItem(storageKey)) return localStorage.getItem(storageKey);
  else
    return window.matchMedia("(prefers-color-scheme: dark)").matches
      ? "dark"
      : "light";
};

const setColorPref = () => {
  localStorage.setItem(storageKey, theme.value);
  reflectPref();
};

const reflectPref = () => {
  document.firstElementChild.setAttribute("data-theme", theme.value);

  document.querySelector("#set-theme")?.setAttribute("aria-label", theme.value);
};

const theme = {
  value: getColorPref(),
};

reflectPref();

window.onload = () => {
  reflectPref();

  document.querySelector("#set-theme").addEventListener("click", onClick);
};

const onClick = () => {
  theme.value = theme.value === "light" ? "dark" : "light";
  setColorPref();
};

window
  .matchMedia("(prefers-color-scheme: dark)")
  .addEventListener("change", ({ matches: isDark }) => {
    theme.value = isDark ? "dark" : "light";
    setColorPref();
  });
