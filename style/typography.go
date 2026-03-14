package style

import "github.com/JohnPitter/openscribe/common"

// Pre-built typography schemes

func ClassicTypography() Typography {
	return Typography{
		HeadingFont: common.NewFont("Georgia", 24).Bold(),
		BodyFont:    common.NewFont("Times New Roman", 12),
		CaptionFont: common.NewFont("Georgia", 10).WithColor(common.Gray),
		CodeFont:    common.NewFont("Courier New", 10),
	}
}

func ModernTypography() Typography {
	return Typography{
		HeadingFont: common.NewFont("Helvetica", 28).Bold(),
		BodyFont:    common.NewFont("Helvetica", 11),
		CaptionFont: common.NewFont("Helvetica", 9).WithColor(common.Gray),
		CodeFont:    common.NewFont("Consolas", 10),
	}
}

func ElegantTypography() Typography {
	return Typography{
		HeadingFont: common.NewFont("Garamond", 26).WithWeight(common.FontWeightSemiBold),
		BodyFont:    common.NewFont("Garamond", 12),
		CaptionFont: common.NewFont("Garamond", 10).Italic().WithColor(common.Gray),
		CodeFont:    common.NewFont("Courier New", 10),
	}
}

func MinimalTypography() Typography {
	return Typography{
		HeadingFont: common.NewFont("Arial", 22).WithWeight(common.FontWeightLight),
		BodyFont:    common.NewFont("Arial", 11),
		CaptionFont: common.NewFont("Arial", 9).WithColor(common.LightGray),
		CodeFont:    common.NewFont("Consolas", 10),
	}
}

func TechTypography() Typography {
	return Typography{
		HeadingFont: common.NewFont("Segoe UI", 26).Bold(),
		BodyFont:    common.NewFont("Segoe UI", 11),
		CaptionFont: common.NewFont("Segoe UI", 9).WithColor(common.Gray),
		CodeFont:    common.NewFont("Cascadia Code", 10),
	}
}
