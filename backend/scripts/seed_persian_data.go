package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()

	// Connect to database
	databaseURL := "postgres://postgres:1@localhost:5432/project_management?sslmode=disable"
	db, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}
	defer db.Close()

	// Delete existing tasks and projects
	log.Println("Cleaning up existing data...")
	_, err = db.Exec(ctx, `DELETE FROM tasks`)
	if err != nil {
		log.Printf("Warning: Failed to delete tasks: %v\n", err)
	} else {
		log.Println("Deleted existing tasks")
	}

	_, err = db.Exec(ctx, `DELETE FROM projects`)
	if err != nil {
		log.Printf("Warning: Failed to delete projects: %v\n", err)
	} else {
		log.Println("Deleted existing projects")
	}

	var adminID *uuid.UUID
	var adminEmail string
	err = db.QueryRow(ctx, `SELECT id, email FROM users WHERE role = 'admin' LIMIT 1`).Scan(&adminID, &adminEmail)
	if err != nil || adminID == nil {
		log.Println("Warning: No admin user found, creating project without owner")
	} else {
		log.Printf("Found admin user: %s (%s)\n", adminEmail, *adminID)
	}

	// Create project with Persian data
	projectID := uuid.New()
	projectTitle := "پروژه تستی توسعه وب‌سایت"
	projectDescription := "این یک پروژه تستی برای آزمایش سیستم مدیریت پروژه است. این پروژه شامل ۱۰۰ تسک با اولویت‌های مختلف و وضعیت‌های گوناگون می‌باشد."
	projectIdentifier := fmt.Sprintf("test-web-dev-%d", time.Now().Unix())
	now := time.Now()

	_, err = db.Exec(ctx, `
		INSERT INTO projects (id, title, description, status, identifier, homepage, is_public, user_id, created_by, created_at, updated_at)
		VALUES ($1, $2, $3, 'active', $4, $5, true, $6, $7, $8, $9)
		ON CONFLICT (id) DO NOTHING
	`, projectID, projectTitle, projectDescription, projectIdentifier, "https://example.com", adminID, adminID, now, now)

	if err != nil {
		log.Fatalf("Failed to create project: %v", err)
	}

	log.Printf("Created project: %s (%s)\n", projectTitle, projectID)

	// Persian task data with various categories, priorities, and statuses
	taskData := []struct {
		title          string
		description    string
		priority       string
		completed      bool
		category       string
		estimatedHours float64
		doneRatio      int
	}{
		// Planning & Design (Tasks 1-10)
		{"تحلیل نیازمندی‌ها", "جمع‌آوری و مستندسازی نیازمندی‌های کسب‌وکار", "High", true, "طراحی", 8.0, 100},
		{"طراحی رابط کاربری", "طراحی وایرفریم و پروتوتایپ رابط کاربری", "High", true, "طراحی", 12.0, 100},
		{"طراحی پایگاه داده", "طراحی شمای پایگاه داده و روابط بین جداول", "High", true, "دیتابیس", 6.0, 100},
		{"تدلیل سیستم", "تدلیل معماری سیستم و انتخاب تکنولوژی‌ها", "High", true, "تحلیل", 10.0, 100},
		{"طراحی API", "طراحی نقاط پایانی API و مستندات", "Medium", true, "طراحی", 8.0, 100},
		{"بررسی امنیت", "بررسی نیازمندی‌های امنیتی و تدابیر لازم", "High", false, "امنیت", 4.0, 50},
		{"طراحی عملکرد", "بررسی و طراحی نیازمندی‌های عملکردی", "Medium", false, "طراحی", 6.0, 30},
		{"بررسی سازگاری", "بررسی سازگاری با مرورگرهای مختلف", "Low", false, "طراحی", 4.0, 20},
		{"طراحی ریسپانسیو", "طراحی نسخه موبایل و تبلت", "Medium", false, "طراحی", 10.0, 40},
		{"نقد و بررسی طرح", "جلسه نقد و بررسی طرح نهایی با تیم", "Medium", false, "ارتباطات", 2.0, 0},

		// Backend Development (Tasks 11-25)
		{"راه‌اندازی پروژه", "ایجاد ساختار پروژه و تنظیمات اولیه", "High", true, "بک‌اند", 2.0, 100},
		{"تنظیم پایگاه داده", "اتصال به PostgreSQL و تنظیم اتصال", "High", true, "دیتابیس", 3.0, 100},
		{"مدل کاربر", "ایجاد مدل کاربر و اعتبارسنجی", "High", true, "بک‌اند", 4.0, 100},
		{"احراز هویت", "پیاده‌سازی سیستم ورود و ثبت‌نام", "High", true, "بک‌اند", 8.0, 100},
		{"مدل پروژه", "ایجاد مدل پروژه و روابط", "Medium", true, "بک‌اند", 3.0, 100},
		{"مدل تسک", "ایجاد مدل تسک و روابط", "Medium", true, "بک‌اند", 3.0, 100},
		{"ریپازیتوری‌ها", "ایجاد لایه دسترسی به داده", "Medium", true, "بک‌اند", 6.0, 100},
		{"سرویس‌ها", "ایجاد لایه سرویس و منطق تجاری", "Medium", false, "بک‌اند", 8.0, 60},
		{"هندلرها", "ایجاد هندلرهای HTTP", "Medium", false, "بک‌اند", 6.0, 50},
		{"مسیریابی", "تعریف مسیرهای API", "Medium", false, "بک‌اند", 2.0, 40},
		{"مدیریت خطاها", "پیاده‌سازی مدیریت خطاهای سراسری", "Medium", false, "بک‌اند", 4.0, 30},
		{"لوگ کردن", "پیاده‌سازی سیستم لاگ‌گیری", "Low", false, "بک‌اند", 2.0, 20},
		{"تست‌های واحد", "نوشتن تست‌های واحد برای سرویس‌ها", "Medium", false, "تست", 10.0, 10},
		{"مستندسازی API", "ایجاد مستندات کامل API", "Low", false, "مستندات", 6.0, 0},
		{"بهینه‌سازی", "بهینه‌سازی کوئری‌های پایگاه داده", "Medium", false, "بک‌اند", 8.0, 0},

		// Frontend Development (Tasks 26-45)
		{"راه‌اندازی فرانت‌اند", "ایجاد پروژه Vite و تنظیمات", "High", true, "فرانت‌اند", 2.0, 100},
		{"کامپوننت‌های پایه", "ایجاد کامپوننت‌های پایه UI", "Medium", true, "فرانت‌اند", 4.0, 100},
		{"صفحه ورود", "پیاده‌سازی فرم ورود", "High", true, "فرانت‌اند", 4.0, 100},
		{"صفحه ثبت‌نام", "پیاده‌سازی فرم ثبت‌نام", "High", true, "فرانت‌اند", 4.0, 100},
		{"داشبورد", "پیاده‌سازی داشبورد اصلی", "High", false, "فرانت‌اند", 8.0, 70},
		{"لیست پروژه‌ها", "پیاده‌سازی نمایش لیست پروژه‌ها", "Medium", false, "فرانت‌اند", 4.0, 60},
		{"فرم پروژه", "پیاده‌سازی فرم ایجاد پروژه", "Medium", false, "فرانت‌اند", 4.0, 50},
		{"لیست تسک‌ها", "پیاده‌سازی نمایش لیست تسک‌ها", "Medium", false, "فرانت‌اند", 6.0, 50},
		{"جزئیات تسک", "پیاده‌سازی صفحه جزئیات تسک", "Medium", false, "فرانت‌اند", 4.0, 40},
		{"فرم تسک", "پیاده‌سازی فرم ایجاد تسک", "Medium", false, "فرانت‌اند", 4.0, 40},
		{"وضعیت تسک", "پیاده‌سازی تغییر وضعیت تسک", "Medium", false, "فرانت‌اند", 2.0, 30},
		{"زمان‌بندی", "پیاده‌سازی ثبت زمان کار", "Medium", false, "فرانت‌اند", 4.0, 30},
		{"فیلترها", "پیاده‌سازی فیلترهای جستجو", "Low", false, "فرانت‌اند", 4.0, 20},
		{"مرتب‌سازی", "پیاده‌سازی مرتب‌سازی تسک‌ها", "Low", false, "فرانت‌اند", 2.0, 20},
		{"پگینیشن", "پیاده‌سازی صفحه‌بندی", "Low", false, "فرانت‌اند", 3.0, 10},
		{"موبایل", "طراحی ریسپانسیو برای موبایل", "Medium", false, "فرانت‌اند", 8.0, 10},
		{"انیمیشن‌ها", "افزودن انیمیشن‌های رابط کاربری", "Low", false, "فرانت‌اند", 6.0, 0},
		{"اعتبارسنجی فرم‌ها", "پیاده‌سازی اعتبارسنجی", "Medium", false, "فرانت‌اند", 4.0, 0},
		{"هشدارها", "پیاده‌سازی سیستم اطلاع‌رسانی", "Medium", false, "فرانت‌اند", 3.0, 0},
		{"لودینگ", "پیاده‌سازی نشانگر لودینگ", "Low", false, "فرانت‌اند", 2.0, 0},

		// Testing (Tasks 46-60)
		{"تست API", "نوشتن تست‌های یکپارچگی API", "High", false, "تست", 12.0, 50},
		{"تست UI", "نوشتن تست‌های رابط کاربری", "Medium", false, "تست", 10.0, 40},
		{"تست لود", "تست بارگذاری سیستم", "Medium", false, "تست", 8.0, 30},
		{"تست امنیت", "تست‌های امنیتی نفوذ", "High", false, "تست", 10.0, 20},
		{"تست مرورگر", "تست در مرورگرهای مختلف", "Medium", false, "تست", 6.0, 20},
		{"تست موبایل", "تست روی دستگاه‌های موبایل", "Medium", false, "تست", 6.0, 10},
		{"تست ریسپانسیو", "تست ریسپانسیو در سایزهای مختلف", "Medium", false, "تست", 4.0, 10},
		{"رفع باگ‌ها", "رفع باگ‌های یافت شده", "High", false, "دیباگ", 20.0, 0},
		{"رفع مشکلات کارایی", "رفع مشکلات کارایی", "High", false, "دیباگ", 12.0, 0},
		{"رفع مشکلات امنیتی", "رفع آسیب‌پذیری‌ها", "High", false, "امنیت", 10.0, 0},
		{"تست رگرسیون", "اجرای تست‌های رگرسیون", "Medium", false, "تست", 8.0, 0},
		{"تست کاربران", "تست با کاربران واقعی", "Medium", false, "تست", 6.0, 0},
		{"جمع‌آوری بازخورد", "جمع‌آوری بازخورد کاربران", "Low", false, "ارتباطات", 4.0, 0},
		{"بهبود تجربه کاربری", "بهبود بر اساس بازخورد", "Medium", false, "طراحی", 8.0, 0},
		{"تست نهایی", "تست نهایی قبل از انتشار", "High", false, "تست", 6.0, 0},

		// Deployment (Tasks 61-70)
		{"تنظیم سرور", "راه‌اندازی سرور پروduction", "High", false, "دیپلوی", 4.0, 60},
		{"کانفیگ Nginx", "تنظیم Nginx برای پروکسی", "High", false, "دیپلوی", 3.0, 50},
		{"کانفیگ SSL", "تنظیم گواهی SSL", "High", false, "دیپلوی", 2.0, 50},
		{"اتصال دیتابیس", "اتصال به دیتابیس پروداکشن", "High", false, "دیپلوی", 2.0, 40},
		{"پیکربندی محیط", "تنظیم متغیرهای محیطی", "Medium", false, "دیپلوی", 2.0, 40},
		{"CI/CD", "تنظیم خط تولید CI/CD", "Medium", false, "دیپلوی", 8.0, 30},
		{"بیلد خودکار", "پیاده‌سازی بیلد خودکار", "Medium", false, "دیپلوی", 4.0, 30},
		{"تست خودکار", "پیاده‌سازی تست‌های خودکار", "Medium", false, "دیپلوی", 4.0, 20},
		{"دیپلوی بک‌اند", "دیپلوی سرور بک‌اند", "High", false, "دیپلوی", 2.0, 0},
		{"دیپلوی فرانت‌اند", "دیپلوی فرانت‌اند", "High", false, "دیپلوی", 2.0, 0},

		// Documentation & Maintenance (Tasks 71-80)
		{"مستندات کاربر", "نوشتن راهنمای کاربر", "Medium", false, "مستندات", 8.0, 20},
		{"مستندات فنی", "نوشتن مستندات فنی", "Medium", false, "مستندات", 10.0, 10},
		{"ویدیوی آموزشی", "ساخت ویدیوهای آموزشی", "Low", false, "مستندات", 12.0, 0},
		{"FAQ", "ایجاد سوالات متداول", "Low", false, "مستندات", 4.0, 0},
		{"مستندات API", "تکمیل مستندات API", "Medium", false, "مستندات", 6.0, 0},
		{"راهنمای نصب", "نوشتن راهنمای نصب و راه‌اندازی", "Medium", false, "مستندات", 4.0, 0},
		{"مانیتورینگ", "راه‌اندازی سیستم مانیتورینگ", "High", false, "دیپلوی", 6.0, 40},
		{"لاگ سرور", "راه‌اندازی لاگ سرور", "Medium", false, "دیپلوی", 4.0, 30},
		{"هشدارها", "تنظیم هشدارهای سیستم", "Medium", false, "دیپلوی", 4.0, 20},
		{"پشتیبان‌گیری", "تنظیم سیستم بکاپ خودکار", "High", false, "دیپلوی", 4.0, 10},

		// Additional Features (Tasks 81-100)
		{"جستجوی پیشرفته", "پیاده‌سازی جستجوی پیشرفته", "Medium", false, "ویژگی", 10.0, 0},
		{"صادرات داده‌ها", "پیاده‌سازی صادرات به Excel", "Medium", false, "ویژگی", 6.0, 0},
		{"وارد کردن داده‌ها", "پیاده‌سازی وارد کردن از CSV", "Medium", false, "ویژگی", 8.0, 0},
		{"نمودارها", "افزودن نمودارهای تحلیلی", "Low", false, "ویژگی", 12.0, 0},
		{"گزارش‌گیری", "پیاده‌سازی سیستم گزارش‌گیری", "Medium", false, "ویژگی", 10.0, 0},
		{"وظایف تکرارشونده", "پیاده‌سازی تسک‌های تکرارشونده", "Low", false, "ویژگی", 8.0, 0},
		{"برچسب‌گذاری", "پیاده‌سازی سیستم برچسب‌ها", "Low", false, "ویژگی", 4.0, 0},
		{"اشتراک‌گذاری", "پیاده‌سازی اشتراک‌گذاری پروژه", "Medium", false, "ویژگی", 6.0, 0},
		{"نظرات", "پیاده‌سازی سیستم نظرات", "Medium", false, "ویژگی", 6.0, 0},
		{"فایل‌های ضمیمه", "پیاده‌سازی آپلود فایل", "Medium", false, "ویژگی", 8.0, 0},
		{"اعلان‌های ایمیلی", "پیاده‌سازی اعلان‌های ایمیلی", "Medium", false, "ویژگی", 8.0, 0},
		{"اعلان‌های مرورگر", "پیاده‌سازی اعلان‌های مرورگر", "Low", false, "ویژگی", 4.0, 0},
		{"تم تاریک", "پیاده‌سازی تم تاریک", "Low", false, "ویژگی", 4.0, 0},
		{"چند زبانه", "پیاده‌سازی چندزبانه", "Medium", false, "ویژگی", 12.0, 0},
		{"تنظیمات کاربر", "پیاده‌سازی تنظیمات کاربر", "Medium", false, "ویژگی", 6.0, 0},
		{"تاریخچه تغییرات", "پیاده‌سازی تاریخچه تغییرات", "Low", false, "ویژگی", 6.0, 0},
		{"حذف نرم", "پیاده‌سازی حذف نرم", "Medium", false, "ویژگی", 4.0, 0},
		{"آرشیو خودکار", "پیاده‌سازی آرشیو خودکار", "Low", false, "ویژگی", 4.0, 0},
		{"آمارهای کاربر", "نمایش آمارهای کاربر", "Low", false, "ویژگی", 6.0, 0},
		{"امتیازدهی", "پیاده‌سازی سیستم امتیازدهی", "Low", false, "ویژگی", 4.0, 0},
	}

	// Insert tasks
	taskCount := 0
	for _, task := range taskData {
		taskID := uuid.New()
		now := time.Now()

		_, err = db.Exec(ctx, `
			INSERT INTO tasks (id, project_id, title, description, priority, completed, category,
				estimated_hours, done_ratio, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			ON CONFLICT (id) DO NOTHING
		`, taskID, projectID, task.title, task.description, task.priority, task.completed,
			task.category, task.estimatedHours, task.doneRatio, now, now)

		if err != nil {
			log.Printf("Failed to insert task %s: %v", task.title, err)
			continue
		}

		taskCount++
	}

	log.Printf("Successfully created %d tasks for project\n", taskCount)
	log.Println("Data population completed!")
}
