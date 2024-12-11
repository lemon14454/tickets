enum Path {
    Login = "Login",
    Register = "Register",
    Main = "Overview"
}

type Events = {
    title: string;
    start_at: string;
    created_at: string;
}

const eventData: Array<Events> = [
    { title: "Justin Bieber", start_at: "2024-11-01 09:00:00", created_at: "2024-10-01 00:00:00" },
    { title: "Justin Bieber", start_at: "2024-11-01 09:00:00", created_at: "2024-10-01 00:00:00" },
    { title: "Justin Bieber", start_at: "2024-11-01 09:00:00", created_at: "2024-10-01 00:00:00" },
    { title: "Justin Bieber", start_at: "2024-11-01 09:00:00", created_at: "2024-10-01 00:00:00" }
]

interface IElementProps {
    type: string;
    text?: string;
    class?: string;
}

function createElement(props: IElementProps): HTMLElement {
    const ele = document.createElement(props.type);
    if (props.text) ele.textContent = props.text;
    if (props.class) ele.classList.add(props.class);
    return ele;
}

interface IFormField {
    text: string;
    type: string;
}

function createForm(fields: Array<IFormField>): HTMLElement {
    const form = document.createElement("form");

    fields.forEach(({ text, type }) => {
        const formGroup = document.createElement("div");
        const label = document.createElement("label");
        label.htmlFor = text;
        label.textContent = text;
        const input = document.createElement("input");
        input.id = text;
        input.type = type;
        formGroup.classList.add("formGroup");
        formGroup.appendChild(label);
        formGroup.appendChild(input);
        form.appendChild(formGroup);
    })

    return form;
}

function renderRegisterMenu(): HTMLElement {
    const form = createForm([
        { text: "Name", type: "string" },
        { text: "Email", type: "email" },
        { text: "Password", type: "password" },
        { text: "Password2", type: "password" },
    ]);

    const btnGroup = createElement({ type: "div", class: "flex" });

    const backBtn = createElement({ type: "button", text: "back" });
    backBtn.addEventListener("click", (e: SubmitEvent) => {
        e.preventDefault();
        render(Path.Login);
    })
    const submitBtn = createElement({ type: "button", text: "Enter" });
    submitBtn.addEventListener("click", (e: SubmitEvent) => {
        // TODO:: Add Login Logic, Redirect To Main Menu
        e.preventDefault();
        render(Path.Login);
    })

    btnGroup.appendChild(backBtn);
    btnGroup.appendChild(submitBtn);
    btnGroup.style.justifyContent = "flex-end";

    form.appendChild(btnGroup);
    return form;
}

function renderLoginMenu(): HTMLElement {
    const form = createForm([
        { text: "Email", type: "email" },
        { text: "Password", type: "password" },
    ]);

    const btnGroup = document.createElement("div");
    btnGroup.classList.add("flex");

    const loginBtn = createElement({ type: "button", text: "Enter" });
    loginBtn.addEventListener("click", (e: SubmitEvent) => {
        e.preventDefault();
        const email = (form.querySelector("#Email") as HTMLInputElement).value;
        const pwd = (form.querySelector("#Password") as HTMLInputElement).value;
        // TODO:: Add Login Logic, Redirect To Main Menu
        render(Path.Main);
    })
    const registerBtn = createElement({ type: "button", text: "Register" });
    registerBtn.addEventListener("click", (e: SubmitEvent) => {
        e.preventDefault();
        // TODO:: Redirect To Register Menu
        render(Path.Register);
    })

    btnGroup.appendChild(registerBtn);
    btnGroup.appendChild(loginBtn);
    btnGroup.style.justifyContent = "flex-end";

    form.appendChild(btnGroup);
    return form;
}

function getLoginUser(): string | null {
    return "User";
}

function renderEventDetail(event: Events, owned: boolean): HTMLElement {
    const wrapper = createElement({ type: "div" });
    const title = createElement({ type: "h3", text: event.title, class: "title" });
    const startAt = createElement({ type: "p", text: `@${event.start_at}`, class: "event_helper" });

    wrapper.appendChild(title);
    wrapper.appendChild(startAt);

    const flex = createElement({ type: "div", class: "flex" });
    flex.style.justifyContent = "flex-end";
    const btn = createElement({ type: "button", text: owned ? "Cancel" : "Buy" });
    btn.addEventListener("click", (e: SubmitEvent) => {
        if (owned) {

        } else {

        }
    })
    flex.appendChild(btn);
    wrapper.appendChild(flex);

    return wrapper;
}

function renderEvent(event: Events, owned: boolean): HTMLElement {
    const wrapper = createElement({ type: "div", class: "event_wrapper" });
    const left = createElement({ type: "div" });
    const t = createElement({ type: "h5", text: event.title, class: "event_title" });
    const startAt = createElement({ type: "p", text: `@${event.start_at}`, class: "event_helper" });

    left.appendChild(t);
    left.appendChild(startAt);

    const detailBtn = createElement({ type: "button", text: ">", class: "event_detail_btn" });
    detailBtn.addEventListener("click", (e: SubmitEvent) => {
        e.preventDefault();
        const eventDetail = renderEventDetail(event, owned);
        renderModal(eventDetail);
    })

    wrapper.appendChild(left);
    wrapper.appendChild(detailBtn);

    return wrapper;
}

function renderEvents(events: Array<Events>, owned: boolean): HTMLElement {
    const container = createElement({ type: "div", class: "event_container" });
    events.forEach((evt) => {
        const event = renderEvent(evt, owned);
        container.appendChild(event);
    })
    return container;
}

function renderMainMenu(): HTMLElement {
    const user = getLoginUser();
    if (!user) {
        render(Path.Login);
        return;
    }

    const wrapper = document.createElement("div");

    const logoutBtn = createElement({ type: "button", text: "Logout", class: "exit_btn" });
    logoutBtn.addEventListener("click", () => {
        render(Path.Login);
    })
    wrapper.appendChild(logoutBtn);

    const greeting = createElement({ type: "p", text: `Welcome ${user}` })
    wrapper.appendChild(greeting);

    const hr = document.createElement("hr");
    wrapper.appendChild(hr);

    const myEventTitle = createElement({ type: "h4", text: "My Events", class: "subtitle" });
    wrapper.appendChild(myEventTitle);

    const myEvents = renderEvents(eventData, true);
    wrapper.appendChild(myEvents);

    const availableEventTitle = createElement({ type: "h4", text: "Available Events", class: "subtitle" });
    wrapper.appendChild(availableEventTitle);

    const availableEvents = renderEvents(eventData, false);
    wrapper.appendChild(availableEvents);

    return wrapper;
}

function renderModal(content: HTMLElement) {
    const root = document.getElementById("root");
    const modal = createElement({ type: "div", class: "modal" });
    const panel = createElement({ type: "div", class: "panel" });
    panel.appendChild(content);

    const closeBtn = createElement({ type: "button", text: "X", class: "exit_btn" });
    closeBtn.addEventListener("click", (e: SubmitEvent) => {
        e.preventDefault();
        root.removeChild(modal);
    })

    panel.appendChild(closeBtn);
    modal.appendChild(panel);
    root.appendChild(modal);
}

function render(path: Path) {
    const root = document.getElementById("root");

    const panel = createElement({ type: "div", class: "panel" });

    /* Title */
    const title = createElement({ type: "h3", text: path, class: "title" });

    /* Content */
    // unmount
    while (root.firstChild) {
        root.removeChild(root.firstChild);
    }

    // mount
    let content: HTMLElement;
    switch (path) {
        case Path.Main:
            content = renderMainMenu();
            break;
        case Path.Login:
            content = renderLoginMenu();
            break;
        case Path.Register:
            content = renderRegisterMenu();
            break;
    }
    panel.appendChild(title);
    panel.appendChild(content);

    /* Append to Root */
    root.appendChild(panel);
}

render(Path.Main);
