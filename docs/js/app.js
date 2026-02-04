console.log("Fetching readme.md...");

async function getReadme() {
    let response = await fetch("https://raw.githubusercontent.com/Kicked-Out/KeyFlip/main/README.md");
    let data = await response.text();

    return data;
}

async function splitText(data) {
    let textList = data.split("\n");
    textList = textList.filter((line) => line.trim() != "");

    return textList;
}

async function addElementToHTML(elementType, className, textContent) {
    let element = document.createElement(elementType);
    element.classList.add(className);
    element.innerHTML = textContent;

    let container = document.querySelector("main");
    container.appendChild(element);

    return element;
}

async function addElementToChildHTML(parentElement, elementType, className, textContent) {
    let element = document.createElement(elementType);

    if (className) {
        element.classList.add(className);
    }
    element.innerHTML = textContent;

    parentElement.appendChild(element);

    return element;
}

async function parseLine(line) {
    let htmlElement = line.replace(/`([^`]+)`/g, "<code class='code-line'>$1</code>");
    htmlElement = htmlElement.replace(/\*\*([^*]+)\*\*/g, "<strong>$1</strong>");

    return htmlElement;
}

async function convertDataToHTML(data) {
    let codeMode = false;
    let codeBlocks = document.querySelectorAll("div.example-block");
    let lastCodeBlock = codeBlocks[codeBlocks.length - 1];

    let tableMode = false;

    let tableBlocks = document.querySelectorAll("table");
    let lastTableBlock;

    if (tableBlocks.length > 0) {
        lastTableBlock = tableBlocks[tableBlocks.length - 1];
    }

    for (let i = 0; i < data.length; i++) {
        let line = data[i];
        let prevLine = data[i - 1] || "";
        let nextLine = data[i + 1] || "";

        switch (true) {
            case line.startsWith(">"):
                tableMode = false;
                line = line.replace("> ", "");
                line = await parseLine(line);

                addElementToHTML("blockquote", "section-quote", line);
                break;
            case line.startsWith("| -"):
                break;
            case line.startsWith("|"):
                if (!tableMode) {
                    tableMode = true;
                    lastTableBlock = await addElementToHTML("table", "section-table", "");

                    let header = await addElementToChildHTML(lastTableBlock, "thead", "", "");
                    let headerRow = await addElementToChildHTML(header, "tr", "", "");
                    let headerCells = line.split("|").filter((cell) => cell.trim() != "");

                    for (let cellText of headerCells) {
                        cellText = await parseLine(cellText.trim());
                        await addElementToChildHTML(headerRow, "th", "", cellText);
                    }
                } else {
                    let body = lastTableBlock.querySelector("tbody");

                    if (!body) {
                        body = await addElementToChildHTML(lastTableBlock, "tbody", "", "");
                    }

                    let bodyRow = await addElementToChildHTML(body, "tr", "", "");
                    let bodyCells = line.split("|").filter((cell) => cell.trim() != "");

                    for (let cellText of bodyCells) {
                        cellText = await parseLine(cellText.trim());
                        await addElementToChildHTML(bodyRow, "td", "", cellText.trim());
                    }

                    if (!nextLine.startsWith("|")) {
                        tableMode = false;
                    }
                }

                break;
            case line.startsWith("```"):
                codeMode = !codeMode;
                tableMode = false;

                if (codeMode) {
                    lastCodeBlock = await addElementToHTML("div", "example-block", "");

                    await addElementToChildHTML(lastCodeBlock, "pre", "example-pre", "");
                }
                break;
            case line.startsWith("---"):
                tableMode = false;
                addElementToHTML("hr", "section-separator", "");
                break;
            case line.startsWith("-"):
                tableMode = false;
                let parsedLine = line.replace("- ", "");
                let listItem = await parseLine(parsedLine);

                if (prevLine.startsWith("-")) {
                    let lists = document.querySelectorAll("ul.section-list");
                    let lastList = lists[lists.length - 1];
                    addElementToChildHTML(lastList, "li", "section-list-item", listItem);
                } else {
                    let ulElement = await addElementToHTML("ul", "section-list", "");
                    await addElementToChildHTML(ulElement, "li", "section-list-item", listItem);
                }
                break;
            case line.startsWith("####"):
                tableMode = false;
                line = line.replace("#### ", "");
                let h4Text = await parseLine(line);
                addElementToHTML("h3", "section-subsection", h4Text);
                break;
            case line.startsWith("###"):
                tableMode = false;
                line = line.replace("### ", "");
                text = await parseLine(line);
                addElementToHTML("h3", "section-subsection", text);
                break;
            case line.startsWith("##"):
                tableMode = false;
                line = line.replace("## ", "");
                let subtitle = await parseLine(line);
                addElementToHTML("h3", "section-subtitle", subtitle);
                break;
            case line.startsWith("#"):
                tableMode = false;
                line = line.replace("# ", "");
                let title = await parseLine(line);
                addElementToHTML("h2", "section-title", title);
                break;
            default:
                tableMode = false;
                line = await parseLine(line);

                if (codeMode) {
                    let preElement = lastCodeBlock.querySelector("pre.example-pre");
                    addElementToChildHTML(preElement, "p", "example-text", line);
                    break;
                }

                addElementToHTML("p", "section-text", line);
                break;
        }
    }
}

async function main() {
    let data = await getReadme();
    let textLines = await splitText(data);

    await convertDataToHTML(textLines);
}

main();
