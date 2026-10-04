import AppKit

let size = 1024
guard let rep = NSBitmapImageRep(bitmapDataPlanes: nil, pixelsWide: size, pixelsHigh: size,
    bitsPerSample: 8, samplesPerPixel: 4, hasAlpha: true, isPlanar: false,
    colorSpaceName: .deviceRGB, bytesPerRow: 0, bitsPerPixel: 0) else { exit(1) }
NSGraphicsContext.saveGraphicsState()
NSGraphicsContext.current = NSGraphicsContext(bitmapImageRep: rep)
let ctx = NSGraphicsContext.current!.cgContext

// 翻转为左上角原点（与常规设计坐标一致）
ctx.translateBy(x: 0, y: CGFloat(size))
ctx.scaleBy(x: 1, y: -1)

func rgb(_ r: Int, _ g: Int, _ b: Int, _ a: CGFloat = 1) -> CGColor {
    CGColor(red: CGFloat(r)/255, green: CGFloat(g)/255, blue: CGFloat(b)/255, alpha: a)
}
func rr(_ x: CGFloat, _ y: CGFloat, _ w: CGFloat, _ h: CGFloat, _ r: CGFloat) -> CGPath {
    CGPath(roundedRect: CGRect(x: x, y: y, width: w, height: h), cornerWidth: r, cornerHeight: r, transform: nil)
}
func fill(_ path: CGPath, _ color: CGColor) {
    ctx.addPath(path); ctx.setFillColor(color); ctx.fillPath()
}

// 背景：蓝色渐变圆角方形（VSCode 风格，占满画布）
ctx.saveGState()
ctx.addPath(rr(0, 0, 1024, 1024, 229))
ctx.clip()
let bgColors = [rgb(10, 108, 196), rgb(53, 162, 245)] as CFArray
let bgGrad = CGGradient(colorsSpace: CGColorSpaceCreateDeviceRGB(), colors: bgColors, locations: [0, 1])!
ctx.drawLinearGradient(bgGrad, start: CGPoint(x: 512, y: 1024), end: CGPoint(x: 512, y: 0), options: [])
ctx.restoreGState()

// 日历主体投影
ctx.saveGState()
ctx.setShadow(offset: CGSize(width: 0, height: 10), blur: 40, color: rgb(0, 0, 20, 0.35))
fill(rr(232, 300, 560, 500, 64), rgb(255, 255, 255))
ctx.restoreGState()

// 顶部两个挂环（白色胶囊 + 内嵌蓝色槽）
for cx: CGFloat in [380, 644] {
    fill(rr(cx - 28, 226, 56, 120, 28), rgb(255, 255, 255))
    fill(rr(cx - 12, 250, 24, 72, 12), rgb(46, 143, 228))
}

// 内部蓝色面板
ctx.saveGState()
ctx.addPath(rr(272, 420, 480, 340, 44))
ctx.clip()
let panelColors = [rgb(22, 104, 200), rgb(47, 155, 239)] as CFArray
let panelGrad = CGGradient(colorsSpace: CGColorSpaceCreateDeviceRGB(), colors: panelColors, locations: [0, 1])!
ctx.drawLinearGradient(panelGrad, start: CGPoint(x: 512, y: 760), end: CGPoint(x: 512, y: 420), options: [])
ctx.restoreGState()

// 面板顶部横条
fill(rr(312, 452, 400, 22, 11), rgb(255, 255, 255))

// 日期网格 5列 x 3行，最后一个为红色高亮
let squareW: CGFloat = 64, squareH: CGFloat = 46
let gapX: CGFloat = 24, gapY: CGFloat = 44
let startX: CGFloat = 320, startY: CGFloat = 500
for row in 0..<3 {
    for col in 0..<5 {
        let x = startX + CGFloat(col) * (squareW + gapX)
        let y = startY + CGFloat(row) * (squareH + gapY)
        let color = (row == 2 && col == 4) ? rgb(229, 72, 77) : rgb(255, 255, 255)
        fill(rr(x, y, squareW, squareH, 10), color)
    }
}

NSGraphicsContext.restoreGraphicsState()

let outURL = URL(fileURLWithPath: CommandLine.arguments[1])
try! rep.representation(using: .png, properties: [:])!.write(to: outURL)
print("written: \(outURL.path)")
