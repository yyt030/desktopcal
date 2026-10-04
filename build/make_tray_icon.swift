import AppKit

// 生成 macOS 菜单栏用的「模板图标」：纯黑不透明像素 + 透明背景。
// 系统会按菜单栏明暗自动反色，所以这里不能使用任何颜色。
//
// 尺寸取 44x44（22pt @2x）。Wails 的托盘实现会把 NSImage 的 size
// 强制设为菜单栏厚度，因此只需保证像素密度足够。
//
// 用法：swift make_tray_icon.swift <输出路径.png>

let size = 44
guard let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: size, pixelsHigh: size,
    bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false,
    colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0) else { exit(1) }
NSGraphicsContext.saveGraphicsState()
NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
let ctx = NSGraphicsContext.current!.cgContext

// 翻转为左上角原点（与设计坐标一致）
ctx.translateBy(x: 0, y: CGFloat(size))
ctx.scaleBy(x: 1, y: -1)

let ink = CGColor(red: 0, green: 0, blue: 0, alpha: 1)
let lineWidth: CGFloat = 2.6

func rr(_ x: CGFloat, _ y: CGFloat, _ w: CGFloat, _ h: CGFloat, _ r: CGFloat) -> CGPath {
    CGPath(roundedRect: CGRect(x: x, y: y, width: w, height: h), cornerWidth: r, cornerHeight: r, transform: nil)
}

func stroke(_ path: CGPath) {
    ctx.addPath(path)
    ctx.setStrokeColor(ink)
    ctx.setLineWidth(lineWidth)
    ctx.strokePath()
}

func fill(_ path: CGPath) {
    ctx.addPath(path)
    ctx.setFillColor(ink)
    ctx.fillPath()
}

// 日历外框
stroke(rr(7, 11, 30, 28, 4.5))

// 顶部两个挂环
fill(rr(13.4, 5, 3.2, 9, 1.6))
fill(rr(27.4, 5, 3.2, 9, 1.6))

// 标题栏分隔线
ctx.move(to: CGPoint(x: 7, y: 19.5))
ctx.addLine(to: CGPoint(x: 37, y: 19.5))
ctx.setStrokeColor(ink)
ctx.setLineWidth(lineWidth)
ctx.strokePath()

// 日期格：2 行 x 3 列
let cell: CGFloat = 4.5
let gap: CGFloat = 4
let startX: CGFloat = 11.25
let startY: CGFloat = 23.5
for row in 0..<2 {
    for col in 0..<3 {
        let x = startX + CGFloat(col) * (cell + gap)
        let y = startY + CGFloat(row) * (cell + gap)
        fill(rr(x, y, cell, cell, 1.2))
    }
}

NSGraphicsContext.restoreGraphicsState()

let outURL = URL(fileURLWithPath: CommandLine.arguments[1])
try! rep.representation(using: .png, properties: [:])!.write(to: outURL)
print("written: \(outURL.path)")
