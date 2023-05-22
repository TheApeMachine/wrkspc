const humanize = () => {
    /*
    mousemove simulates a mousemove event along the given path.
    */
    const movemouse = (path) => {
        for (let i = 1; i < path.length; i++) {
            const event = new MouseEvent('mousemove', { clientX: path[i].x, clientY: path[i].y });
            window.dispatchEvent(event);
        }
    }

    /*
    scrollmouse simulates a scroll event.
    */
    const scrollmouse = (delta) => {
        const event = new WheelEvent('wheel', { deltaY: delta });
        window.dispatchEvent(event);
    }

    /*
    mouseclick simulates a mouseclick event at the current position, with the given button.
    */
    const clickmouse = (button) => {
        const event = new MouseEvent('click', { button: {"left": 0, "right": 1}[button]});
        window.dispatchEvent(event);
    }

    /*
    kepress simulates a keypress event on the given element.
    */
    const presskey = (key) => {
        const event = new KeyboardEvent('keypress', { key: key });
        window.dispatchEvent(event);
    }

    /*
    getpage returns the page at the given url as a string.
    */
    const getpage = () => {
        return document.documentElement.outerHTML;
    }

    return { movemouse, scrollmouse, clickmouse, presskey, getpage }
};
