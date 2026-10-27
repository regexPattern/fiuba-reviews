use actualizador_ofertas::gui;
use gpui_kit::AppContext;

fn main() {
    gpui_kit::application().run(|cx| {
        gpui_kit::init(cx);
        gpui_kit::open_window(gpui_kit::WindowOptions::default(), cx, |_, cx| {
            cx.new(|_| gui::HelloWorld)
        })
        .expect("Failed to open window");
    });
}
