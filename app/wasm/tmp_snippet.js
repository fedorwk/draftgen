function download(filename, text) {
  var element = document.createElement('a');
  element.setAttribute('href', 'data:application/octet-stream,' + encodeURIComponent(text));
  element.setAttribute('download', filename);

  element.style.display = 'none';
  document.body.appendChild(element);

  element.click();

  document.body.removeChild(element);
}


// Document Download
// Допустим, у тебя метод generateSomeShit вызывает васм говнище и вовзращает результат в виде массива, тогда:
const shit = await generateSomeShit();
const blob = new Blob([shit]);
const url = URL.createObjectURL(blob);
const a = document.createElement('a');
a.href = url;
a.downolad = "some-binary.bin"
document.body.append(a);
a.click();
a.remove();
URL.revokeObjectURL(url);